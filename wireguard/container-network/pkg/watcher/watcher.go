// Package watcher implements container event watching and discovery.
package watcher

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"

	"container-network/pkg/client"
)

// ContainerEventType represents the type of container event.
type ContainerEventType int

const (
	// ContainerStarted indicates a container has started.
	ContainerStarted ContainerEventType = iota
	// ContainerStopped indicates a container has stopped (died or crashed)
	ContainerStopped
)

func (t ContainerEventType) String() string {
	switch t {
	case ContainerStarted:
		return "started"
	case ContainerStopped:
		return "stopped"
	default:
		return "unknown"
	}
}

// PortMapping represents a container port mapping.
type PortMapping struct {
	HostIP        string
	HostPort      uint16
	ContainerPort uint16
	Protocol      string
}

// ContainerInfo contains relevant container information for events.
type ContainerInfo struct {
	ID          string
	Name        string
	IPAddress   string
	NetworkName string
	Ports       []PortMapping
	Labels      map[string]string
}

// ContainerEvent represents an event about a container.
type ContainerEvent struct {
	Type      ContainerEventType
	Container ContainerInfo
	Timestamp time.Time
}

// Config contains watcher configuration.
type Config struct {
	NetworkName string
	EnableLabel string
}

// DefaultConfig returns the default watcher configuration.
func DefaultConfig() Config {
	return Config{
		NetworkName: "bridge",
		EnableLabel: "",
	}
}

// Watcher watches for container events and reports them.
type Watcher struct {
	client          *client.Client
	config          Config
	events          chan ContainerEvent
	knownContainers map[string]ContainerInfo
	mu              sync.RWMutex
}

// NewWatcher creates a new container watcher.
func NewWatcher(c *client.Client, config Config) *Watcher {
	return &Watcher{
		client:          c,
		config:          config,
		events:          make(chan ContainerEvent, 200),
		knownContainers: make(map[string]ContainerInfo),
	}
}

// addKnownContainer adds a container to the known containers map.
func (w *Watcher) addKnownContainer(id string, info ContainerInfo) {
	w.mu.Lock()
	w.knownContainers[id] = info
	w.mu.Unlock()
}

// removeKnownContainer removes a container from the known containers map
// and returns the container info if it was previously known.
func (w *Watcher) removeKnownContainer(id string) (ContainerInfo, bool) {
	w.mu.Lock()
	info, wasKnown := w.knownContainers[id]
	delete(w.knownContainers, id)
	w.mu.Unlock()
	return info, wasKnown
}

// Events returns the channel that receives container events.
func (w *Watcher) Events() <-chan ContainerEvent {
	return w.events
}

// Start begins watching for container events.
func (w *Watcher) Start(ctx context.Context) error {
	if err := w.discoverExistingContainers(ctx); err != nil {
		return fmt.Errorf("discovering existing containers: %w", err)
	}
	go w.watchEvents(ctx)
	return nil
}

func (w *Watcher) discoverExistingContainers(ctx context.Context) error {
	filters := map[string][]string{
		"network": {w.config.NetworkName},
	}
	if w.config.EnableLabel != "" {
		filters["label"] = []string{w.config.EnableLabel}
	}
	containers, err := w.client.ListContainers(ctx, filters)
	if err != nil {
		return fmt.Errorf("listing containers: %w", err)
	}
	for _, container := range containers {
		if err := w.processContainer(ctx, &container); err != nil {
			return err
		}
	}
	return err
}

func (w *Watcher) processContainer(ctx context.Context, container *client.Container) error {
	if !w.shouldWatch(container) {
		return nil
	}
	info := w.extractContainerInfo(container)
	w.addKnownContainer(container.ID, info)
	select {
	case w.events <- ContainerEvent{
		Type:      ContainerStarted,
		Container: info,
		Timestamp: time.Now(),
	}:
	case <-ctx.Done():
		return ctx.Err()
	}
	return nil
}

func (w *Watcher) shouldWatch(container *client.Container) bool {
	// Check label filter if configured
	if w.config.EnableLabel != "" {
		labelValue, hasLabel := container.Labels[w.config.EnableLabel]
		if !hasLabel || !strings.EqualFold(labelValue, "true") {
			return false
		}
	}
	// Check network filter
	if container.NetworkSettings != nil && container.NetworkSettings.Networks != nil {
		_, hasNetwork := container.NetworkSettings.Networks[w.config.NetworkName]
		return hasNetwork
	}
	return false
}

func (w *Watcher) extractContainerInfo(container *client.Container) ContainerInfo {
	info := ContainerInfo{
		ID:          container.ID,
		NetworkName: w.config.NetworkName,
		Labels:      container.Labels,
	}
	if len(container.Names) > 0 {
		info.Name = strings.TrimPrefix(container.Names[0], "/")
	}
	if container.NetworkSettings != nil && container.NetworkSettings.Networks != nil {
		if network, ok := container.NetworkSettings.Networks[w.config.NetworkName]; ok {
			info.IPAddress = network.IPAddress
		}
	}
	for _, port := range container.Ports {
		if port.PublicPort > 0 {
			info.Ports = append(info.Ports, PortMapping{
				HostIP:        port.IP,
				HostPort:      port.PublicPort,
				ContainerPort: port.PrivatePort,
				Protocol:      port.Type,
			})
		}
	}
	return info
}

func (w *Watcher) watchEvents(ctx context.Context) {
	filters := map[string][]string{
		"type":  {"container"},
		"event": {"start", "stop", "die", "kill"},
	}
	if w.config.EnableLabel != "" {
		filters["label"] = []string{w.config.EnableLabel}
	}
	eventCh, errCh := w.client.Events(ctx, filters)
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-errCh:
			if err != nil && ctx.Err() == nil {
				slog.Error("Error watching events", "error", err)
				time.Sleep(2 * time.Second)
				go w.watchEvents(ctx)
				return
			}
		case event, ok := <-eventCh:
			if !ok {
				slog.Error("Event channel closed")
				return
			}
			w.handleEvent(ctx, event)
		}
	}
}

func (w *Watcher) handleEvent(ctx context.Context, event client.Event) {
	if event.Type != "container" {
		return
	}
	containerID := event.Actor.ID
	// Docker uses Action, Podman uses Status
	action := event.Action
	if action == "" {
		action = event.Status
	}
	var eventType ContainerEventType
	switch {
	case strings.HasPrefix(action, "start"):
		eventType = ContainerStarted
	case strings.HasPrefix(action, "stop"):
		eventType = ContainerStopped
	case strings.HasPrefix(action, "die"), strings.HasPrefix(action, "kill"):
		eventType = ContainerStopped
	default:
		return
	}
	if eventType == ContainerStopped {
		info, wasKnown := w.removeKnownContainer(containerID)
		if wasKnown {
			select {
			case w.events <- ContainerEvent{
				Type:      eventType,
				Container: info,
				Timestamp: time.Unix(event.Time, 0),
			}:
			case <-ctx.Done():
				return
			}
		}
	} else {
		filters := map[string][]string{
			"id": {containerID},
		}
		containers, err := w.client.ListContainers(ctx, filters)
		if err != nil {
			slog.Error("Error retrieving container", "containerID", containerID, "error", err)
			return
		}
		for _, container := range containers {
			w.processContainer(ctx, &container)
		}
	}
}
