package main

import (
	"bufio"
	"bytes"
	"context"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"

	"container-network/pkg/client"
	"container-network/pkg/config"
	"container-network/pkg/handler"
	"container-network/pkg/watcher"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}
	// Run startup script if configured
	if cfg.StartupScript != "" {
		slog.Info("Running startup script", "script", cfg.StartupScript)
		if err := runScript(cfg.StartupScript); err != nil {
			slog.Error("Startup script failed", "error", err)
			os.Exit(1)
		}
		slog.Info("Startup script completed successfully")
	}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	wg.Add(1)
	go func() {
		defer wg.Done()
		sig := <-sigCh
		slog.Info("Received signal, shutting down...", "signal", sig)
		cancel()
		// Run shutdown script if configured
		if cfg.ShutdownScript != "" {
			slog.Info("Running shutdown script", "script", cfg.ShutdownScript)
			if err := runScript(cfg.ShutdownScript); err != nil {
				slog.Error("Shutdown script failed", "error", err)
			} else {
				slog.Info("Shutdown script completed successfully")
			}
		}
	}()
	dockerClient, err := client.NewClient(cfg.RuntimeAPI)
	if err != nil {
		slog.Error("Failed to create container client", "error", err)
		os.Exit(1)
	}
	slog.Info("Connecting to container runtime", "api", cfg.RuntimeAPI)
	if err := dockerClient.Ping(ctx); err != nil {
		slog.Error("Failed to connect to container runtime", "error", err)
		os.Exit(1)
	}
	slog.Info("Successfully connected to container runtime")
	watcherConfig := watcher.Config{
		NetworkName: cfg.WatchNetwork,
		EnableLabel: cfg.WatchContainerLabel,
	}
	w := watcher.NewWatcher(dockerClient, watcherConfig)
	h := handler.NewHandler(w.Events(), cfg.IptablesMangleMarkPublishedPorts, cfg.IptablesDnatPortsLabel)
	go func() {
		if err := h.Start(ctx); err != nil && err != context.Canceled {
			slog.Error("Event handler error", "error", err)
		}
	}()
	if cfg.WatchContainerLabel != "" {
		slog.Info("Starting container watcher", "network", cfg.WatchNetwork, "label", cfg.WatchContainerLabel)
	} else {
		slog.Info("Starting container watcher", "network", cfg.WatchNetwork)
	}
	if err := w.Start(ctx); err != nil {
		slog.Error("Failed to start watcher", "error", err)
		os.Exit(1)
	}
	slog.Info("Watching for container events. Press Ctrl+C to stop.")
	<-ctx.Done()
	wg.Wait()
	slog.Info("Shutdown complete")
}

// runScript executes the given script and returns any error.
func runScript(script string) error {
	cmd := exec.Command("/bin/sh", "-c", script)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	scanner := bufio.NewScanner(&stdout)
	for scanner.Scan() {
		slog.Info(scanner.Text())
	}
	scanner = bufio.NewScanner(&stderr)
	for scanner.Scan() {
		slog.Error(scanner.Text())
	}
	return err
}
