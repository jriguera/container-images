// Package handler implements event handling for container events.
package handler

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"container-network/pkg/watcher"
)

// Handler processes container events.
type Handler struct {
	events                           <-chan watcher.ContainerEvent
	iptablesMangleMarkPublishedPorts string
	iptablesDnatPortsLabel           string
}

// port represents a port with protocol for iptables rules.
type port struct {
	port     uint16
	protocol string
}

const (
	warmUpMaxAttempts = 60
	warmUpTimeout     = 1 * time.Second
	warmUpInterval    = 1 * time.Second
)

// NewHandler creates a new event handler.
// iptablesMangleMarkPublishedPorts is the iptables mark value to use for published ports.
// If empty, iptables rules are not created.
// iptablesDnatPortsLabel is the label name containing DNAT port mappings.
func NewHandler(events <-chan watcher.ContainerEvent, iptablesMangleMarkPublishedPorts, iptablesDnatPortsLabel string) *Handler {
	return &Handler{
		events:                           events,
		iptablesMangleMarkPublishedPorts: iptablesMangleMarkPublishedPorts,
		iptablesDnatPortsLabel:           iptablesDnatPortsLabel,
	}
}

// Start begins processing events.
func (h *Handler) Start(ctx context.Context) error {
	slog.Info("Event handler started, waiting for container events...")
	for {
		select {
		case <-ctx.Done():
			slog.Info("Event handler stopping...")
			return ctx.Err()
		case event, ok := <-h.events:
			if !ok {
				slog.Info("Event channel closed")
				return nil
			}
			switch event.Type {
			case watcher.ContainerStarted:
				go h.handleContainerStarted(event)
			case watcher.ContainerStopped:
				go h.handleContainerStopped(event)
			}
		}
	}
}

func (h *Handler) handleContainerStarted(event watcher.ContainerEvent) {
	c := event.Container
	logger := slog.With("container", c.Name, "containerID", c.ID[:12], "timestamp", event.Timestamp.Format("2006-01-02 15:04:05"))
	logger.Info("Handling container started")
	if c.IPAddress != "" {
		logger = logger.With("ip", c.IPAddress)
		var cPort uint16
		var cProtocol string
		var dnatPorts []port

		if len(c.Ports) > 0 {
			// Use the first published port for warm-up
			cPort = c.Ports[0].ContainerPort
			cProtocol = c.Ports[0].Protocol
		}
		h.warmupReversePath(logger, c.IPAddress, cPort, cProtocol)
		// Get DNAT ports from label (if any)
		if h.iptablesDnatPortsLabel != "" {
			if portsValue, ok := c.Labels[h.iptablesDnatPortsLabel]; ok {
				dnatPorts = parsePorts(logger, portsValue)
				h.addIptablesDNATRules(logger, c.IPAddress, dnatPorts)
			}
		}
		// Add iptables mangle rules for published ports (excluding DNAT ports)
		if h.iptablesMangleMarkPublishedPorts != "" {
			portsToMark := filterPublishedPorts(c.Ports, dnatPorts)
			h.addIptablesMarkRules(logger, portsToMark)
		}
	}
}

func (h *Handler) handleContainerStopped(event watcher.ContainerEvent) {
	c := event.Container
	logger := slog.With("container", c.Name, "containerID", c.ID[:12], "timestamp", event.Timestamp.Format("2006-01-02 15:04:05"))
	logger.Info("Handling container stopped")
	if c.IPAddress != "" {
		logger = logger.With("ip", c.IPAddress)
		var dnatPorts []port
		// Get DNAT ports from label (if any)
		if h.iptablesDnatPortsLabel != "" {
			if portsValue, ok := c.Labels[h.iptablesDnatPortsLabel]; ok {
				dnatPorts = parsePorts(logger, portsValue)
				h.removeIptablesDNATRules(logger, c.IPAddress, dnatPorts)
			}
		}
		// Remove iptables mangle rules for published ports (excluding DNAT ports)
		if h.iptablesMangleMarkPublishedPorts != "" {
			portsToUnmark := filterPublishedPorts(c.Ports, dnatPorts)
			h.removeIptablesMarkRules(logger, portsToUnmark)
		}
	}
}

// parsePorts parses a comma-separated list of ports in the format "port[/protocol]".
// If no protocol is specified, defaults to "tcp".
// Example: "80,443/tcp,53/udp" -> [{80, "tcp"}, {443, "tcp"}, {53, "udp"}]
func parsePorts(logger *slog.Logger, portsStr string) []port {
	var ports []port
	for _, p := range strings.Split(portsStr, ",") {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		parts := strings.SplitN(p, "/", 2)
		portNum, err := strconv.ParseUint(parts[0], 10, 16)
		if err != nil {
			logger.Warn("Invalid port number", "port", parts[0])
			continue
		}
		protocol := "tcp"
		if len(parts) == 2 {
			protocol = strings.ToLower(parts[1])
		}
		ports = append(ports, port{port: uint16(portNum), protocol: protocol})
	}
	return ports
}

// filterPublishedPorts returns published ports that are not in the DNAT ports list.
// Matching is done by port number and protocol.
func filterPublishedPorts(publishedPorts []watcher.PortMapping, dnatPorts []port) []port {
	var filtered []port
	dnatSet := make(map[uint16]string)
	for _, p := range dnatPorts {
		dnatSet[p.port] = p.protocol
	}
	for _, p := range publishedPorts {
		if protocol, ok := dnatSet[p.HostPort]; ok {
			if protocol == p.Protocol {
				continue
			}
		}
		filtered = append(filtered, port{port: p.HostPort, protocol: p.Protocol})
	}
	return filtered
}

// warmupReversePath pings the container IP until it responds to warm up the Linux reverse path filter routing tables.
// If port is provided, first tries TCP/UDP connection to the specified port. If that fails after all attempts, falls back to ICMP.
// If port is 0, uses ICMP ping directly.
func (h *Handler) warmupReversePath(logger *slog.Logger, ip string, port uint16, protocol string) bool {
	// Try with port first if provided
	if port > 0 {
		portLogger := logger.With("port", port, "protocol", protocol)
		portLogger.Info("Warming up reverse path")
		if h.tryConnect(portLogger, protocol, fmt.Sprintf("%s:%d", ip, port)) {
			return true
		}
		portLogger.Warn("Failed to warm up via port, falling back to ICMP")
	}
	// Try ICMP
	logger.Info("Warming up reverse path (ICMP)")
	return h.tryConnect(logger, "ip4:icmp", ip)
}

// tryConnect attempts to connect to the given address with the given protocol.
func (h *Handler) tryConnect(logger *slog.Logger, protocol, address string) bool {
	for i := 0; i < warmUpMaxAttempts; i++ {
		conn, err := net.DialTimeout(protocol, address, warmUpTimeout)
		if err == nil {
			conn.Close()
			logger.Info("Reverse path warmed up", "attempt", i+1)
			return true
		}
		time.Sleep(warmUpInterval)
	}
	logger.Warn("Failed to warm up reverse path", "maxAttempts", warmUpMaxAttempts)
	return false
}

// addIptablesMarkRules adds iptables mangle PREROUTING rules to mark packets from published ports with mark 2.
func (h *Handler) addIptablesMarkRules(logger *slog.Logger, ports []port) {
	for _, p := range ports {
		if err := h.iptablesMarkRule(logger, "-A", p.protocol, p.port); err != nil {
			logger.Error("Failed to add iptables mark rule", "port", p.port, "protocol", p.protocol, "error", err)
		} else {
			logger.Info("Added iptables mark rule", "port", p.port, "protocol", p.protocol)
		}
	}
}

// removeIptablesMarkRules removes iptables mangle PREROUTING rules for the specified published ports.
func (h *Handler) removeIptablesMarkRules(logger *slog.Logger, ports []port) {
	for _, p := range ports {
		if err := h.iptablesMarkRule(logger, "-D", p.protocol, p.port); err != nil {
			logger.Error("Failed to remove iptables mark rule", "port", p.port, "protocol", p.protocol, "error", err)
		} else {
			logger.Info("Removed iptables mark rule", "port", p.port, "protocol", p.protocol)
		}
	}
}

// iptablesMarkRule executes an iptables command to add (-A) or delete (-D) a
// mangle PREROUTING rule that marks packets with --sport <port>.
func (h *Handler) iptablesMarkRule(logger *slog.Logger, action, protocol string, port uint16) error {
	// iptables -t mangle -A PREROUTING -p <protocol> --sport <port> -j MARK --set-mark <value>
	args := []string{
		"-t", "mangle",
		action, "PREROUTING",
		"-p", protocol,
		"--sport", fmt.Sprintf("%d", port),
		"-j", "MARK",
		"--set-mark", h.iptablesMangleMarkPublishedPorts,
	}
	logger.Debug("Executing iptables", "args", strings.Join(args, " "))
	cmd := exec.Command("iptables", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

// addIptablesDNATRules adds DNAT and FORWARD rules for the specified ports.
func (h *Handler) addIptablesDNATRules(logger *slog.Logger, containerIP string, ports []port) {
	for _, p := range ports {
		if err := h.iptablesDNATRule(logger, "-A", p.protocol, p.port, containerIP); err != nil {
			logger.Error("Failed to add DNAT rule", "port", p.port, "protocol", p.protocol, "error", err)
		} else {
			logger.Info("Added DNAT rule", "port", p.port, "protocol", p.protocol)
		}
		if err := h.iptablesForwardRule(logger, "-A", p.protocol, p.port, containerIP); err != nil {
			logger.Error("Failed to add FORWARD rule", "port", p.port, "protocol", p.protocol, "error", err)
		} else {
			logger.Info("Added FORWARD rule", "port", p.port, "protocol", p.protocol)
		}
	}
}

// removeIptablesDNATRules removes DNAT and FORWARD rules for the specified ports.
func (h *Handler) removeIptablesDNATRules(logger *slog.Logger, containerIP string, ports []port) {
	for _, p := range ports {
		if err := h.iptablesDNATRule(logger, "-D", p.protocol, p.port, containerIP); err != nil {
			logger.Error("Failed to remove DNAT rule", "port", p.port, "protocol", p.protocol, "error", err)
		} else {
			logger.Info("Removed DNAT rule", "port", p.port, "protocol", p.protocol)
		}
		if err := h.iptablesForwardRule(logger, "-D", p.protocol, p.port, containerIP); err != nil {
			logger.Error("Failed to remove FORWARD rule", "port", p.port, "protocol", p.protocol, "error", err)
		} else {
			logger.Info("Removed FORWARD rule", "port", p.port, "protocol", p.protocol)
		}
	}
}

// iptablesDNATRule executes a nat PREROUTING DNAT rule.
func (h *Handler) iptablesDNATRule(logger *slog.Logger, action, protocol string, port uint16, containerIP string) error {
	// iptables -t nat -A PREROUTING -p <protocol> --dport <port> -j DNAT --to-destination <containerip>:<port>
	args := []string{
		"-t", "nat",
		action, "PREROUTING",
		"-p", protocol,
		"--dport", fmt.Sprintf("%d", port),
		"-j", "DNAT",
		"--to-destination", fmt.Sprintf("%s:%d", containerIP, port),
	}
	logger.Debug("Executing iptables", "args", strings.Join(args, " "))
	cmd := exec.Command("iptables", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}

// iptablesForwardRule executes a FORWARD rule to accept traffic to a container.
func (h *Handler) iptablesForwardRule(logger *slog.Logger, action, protocol string, port uint16, containerIP string) error {
	// iptables -A FORWARD -p <protocol> -d <containerip> --dport <port> -j ACCEPT
	args := []string{
		action, "FORWARD",
		"-p", protocol,
		"-d", containerIP,
		"--dport", fmt.Sprintf("%d", port),
		"-j", "ACCEPT",
	}
	logger.Debug("Executing iptables", "args", strings.Join(args, " "))
	cmd := exec.Command("iptables", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%v: %s", err, string(output))
	}
	return nil
}
