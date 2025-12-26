// Package config provides configuration management with support for
// command-line flags and environment variables.
package config

import (
	"flag"
	"fmt"
	"os"
)

// AppName is the name of the application.
const AppName = "container-network"

// Version is set at build time using -ldflags "-X container-network/pkg/config.Version=x.y.z"
var Version = "dev"

// Config holds the application configuration.
type Config struct {
	RuntimeAPI                       string
	WatchNetwork                     string
	WatchContainerLabel              string
	IptablesMangleMarkPublishedPorts string
	IptablesDnatPortsLabel           string
	StartupScript                    string
	ShutdownScript                   string
}

// Default socket paths for Docker and Podman
const (
	DefaultDockerSocket     = "/var/run/docker.sock"
	DefaultPodmanSocket     = "/run/podman/podman.sock"
	DefaultPodmanUserSocket = "/run/user/%d/podman/podman.sock"
)

// DefaultConfig returns the default configuration.
func DefaultConfig() *Config {
	return &Config{
		RuntimeAPI:             detectDefaultSocket(),
		WatchNetwork:           "bridge",
		WatchContainerLabel:    "network.enable",
		IptablesDnatPortsLabel: "network.dnat.ports",
	}
}

// getStringFlag returns the flag value if set, otherwise the env var value if set,
// otherwise the default value.
func getStringFlag(flagVal *string, envKey, defaultVal string) string {
	if *flagVal != "" {
		return *flagVal
	}
	if env := os.Getenv(envKey); env != "" {
		return env
	}
	return defaultVal
}

func detectDefaultSocket() string {
	if _, err := os.Stat(DefaultDockerSocket); err == nil {
		return DefaultDockerSocket
	}
	if _, err := os.Stat(DefaultPodmanSocket); err == nil {
		return DefaultPodmanSocket
	}
	uid := os.Getuid()
	userSocket := fmt.Sprintf("/run/user/%d/podman/podman.sock", uid)
	if _, err := os.Stat(userSocket); err == nil {
		return userSocket
	}
	homeDir, _ := os.UserHomeDir()
	if homeDir != "" {
		macDockerSocket := homeDir + "/.docker/run/docker.sock"
		if _, err := os.Stat(macDockerSocket); err == nil {
			return macDockerSocket
		}
	}
	return DefaultDockerSocket
}

// Load loads configuration from flags and environment variables.
func Load() (*Config, error) {
	cfg := DefaultConfig()
	runtimeAPI := flag.String("runtime-api", "", "Path to Docker/Podman socket (env: RUNTIME_API, default: auto-detect)")
	watchNetwork := flag.String("watch-network", "", "Network name to watch (env: WATCH_NETWORK, default: bridge)")
	watchContainerLabel := flag.String("watch-container-label", "", "Label name to enable watching (env: WATCH_CONTAINER_LABEL)")
	iptablesMangleMark := flag.String("iptables-mangle-mark-published-ports", "", "iptables mark value for published ports (env: IPTABLES_MANGLE_MARK_PUBLISHED_PORTS)")
	iptablesDnatPortsLabel := flag.String("iptables-dnat-ports-label", "", "Label name for DNAT ports (env: IPTABLES_DNAT_PORTS_LABEL, default: network.dnat.ports)")
	startupScript := flag.String("startup-script", "", "Script to run before starting - exit non-zero to abort (env: STARTUP_SCRIPT)")
	shutdownScript := flag.String("shutdown-script", "", "Script to run before shutdown (env: SHUTDOWN_SCRIPT)")
	showHelp := flag.Bool("help", false, "Show help message")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Usage = printUsage
	flag.Parse()
	if *showHelp {
		printUsage()
		os.Exit(0)
	}
	if *showVersion {
		fmt.Printf("%s %s\n", AppName, Version)
		os.Exit(0)
	}
	cfg.RuntimeAPI = getStringFlag(runtimeAPI, "RUNTIME_API", cfg.RuntimeAPI)
	cfg.WatchNetwork = getStringFlag(watchNetwork, "WATCH_NETWORK", cfg.WatchNetwork)
	cfg.WatchContainerLabel = getStringFlag(watchContainerLabel, "WATCH_CONTAINER_LABEL", cfg.WatchContainerLabel)
	cfg.IptablesMangleMarkPublishedPorts = getStringFlag(iptablesMangleMark, "IPTABLES_MANGLE_MARK_PUBLISHED_PORTS", cfg.IptablesMangleMarkPublishedPorts)
	cfg.IptablesDnatPortsLabel = getStringFlag(iptablesDnatPortsLabel, "IPTABLES_DNAT_PORTS_LABEL", cfg.IptablesDnatPortsLabel)
	cfg.StartupScript = getStringFlag(startupScript, "STARTUP_SCRIPT", cfg.StartupScript)
	cfg.ShutdownScript = getStringFlag(shutdownScript, "SHUTDOWN_SCRIPT", cfg.ShutdownScript)
	return cfg, nil
}

func printUsage() {
	w := flag.CommandLine.Output()
	fmt.Fprintf(w, "%s - Watch Docker/Podman containers on a network\n\n", AppName)
	fmt.Fprintf(w, "Usage: %s [OPTIONS]\n\n", AppName)
	fmt.Fprintln(w, "Options:")
	flag.PrintDefaults()
	fmt.Fprintf(w, `
Description:
  This daemon watches for Docker or Podman containers that are attached to a
  specific network and optionally have a specific label set to "true".

  When a matching container starts, it can execute iptables rules for DNAT
  and/or mark published ports. When a container stops, the rules are removed.

Examples:
  # Watch containers on the default bridge network
  %[1]s

  # Watch containers on a custom network with a label filter
  %[1]s -watch-network my-network -watch-container-label network.rp.enable

  # Use Podman socket explicitly
  %[1]s -runtime-api /run/podman/podman.sock

  # Use environment variables
  WATCH_NETWORK=my-network %[1]s
`, AppName)
}
