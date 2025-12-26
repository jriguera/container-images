// Package client provides a minimal Docker/Podman API client
// that works with both container runtimes via Unix socket or HTTP.
package client

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Container represents a container with its relevant information.
type Container struct {
	ID              string            `json:"Id"`
	Names           []string          `json:"Names"`
	Labels          map[string]string `json:"Labels"`
	State           string            `json:"State"`
	Status          string            `json:"Status"`
	NetworkSettings *NetworkSettings  `json:"NetworkSettings"`
	Ports           []Port            `json:"Ports"`
}

// PortBinding represents a host port binding.
type PortBinding struct {
	HostIP   string `json:"HostIp"`
	HostPort string `json:"HostPort"`
}

// NetworkSettings contains container network configuration.
type NetworkSettings struct {
	Networks map[string]*NetworkEndpoint `json:"Networks"`
	Ports    map[string][]PortBinding    `json:"Ports"`
}

// NetworkEndpoint represents a container's endpoint in a network.
type NetworkEndpoint struct {
	IPAddress   string `json:"IPAddress"`
	NetworkID   string `json:"NetworkID"`
	Gateway     string `json:"Gateway"`
	MacAddress  string `json:"MacAddress"`
	IPPrefixLen int    `json:"IPPrefixLen"`
}

// Port represents a port mapping.
type Port struct {
	IP          string `json:"IP,omitempty"`
	PrivatePort uint16 `json:"PrivatePort"`
	PublicPort  uint16 `json:"PublicPort,omitempty"`
	Type        string `json:"Type"`
}

// Event represents a container event from the Docker/Podman API.
type Event struct {
	Type   string `json:"Type"`
	Action string `json:"Action"`
	Actor  Actor  `json:"Actor"`
	Time   int64  `json:"time"`
	Status string `json:"status"`
}

// Actor contains information about the object that triggered the event.
type Actor struct {
	ID         string            `json:"ID"`
	Attributes map[string]string `json:"Attributes"`
}

// Client is a minimal Docker/Podman API client.
type Client struct {
	httpClient *http.Client
	host       string
	scheme     string
	basePath   string
}

// DefaultAPIVersion is the Docker/Podman API version used by the client.
// Version 1.41 corresponds to Docker 20.10.x and is compatible with Podman 3.x/4.x.
const DefaultAPIVersion = "1.41"

// ContainerInspect contains detailed container information.
type ContainerInspect struct {
	ID              string          `json:"Id"`
	Name            string          `json:"Name"`
	State           ContainerState  `json:"State"`
	Config          ContainerConfig `json:"Config"`
	NetworkSettings NetworkSettings `json:"NetworkSettings"`
}

// ContainerState represents the state of a container.
type ContainerState struct {
	Status     string `json:"Status"`
	Running    bool   `json:"Running"`
	Paused     bool   `json:"Paused"`
	Restarting bool   `json:"Restarting"`
	OOMKilled  bool   `json:"OOMKilled"`
	Dead       bool   `json:"Dead"`
	Pid        int    `json:"Pid"`
	ExitCode   int    `json:"ExitCode"`
	StartedAt  string `json:"StartedAt"`
	FinishedAt string `json:"FinishedAt"`
}

// ContainerConfig contains container configuration.
type ContainerConfig struct {
	Labels       map[string]string   `json:"Labels"`
	ExposedPorts map[string]struct{} `json:"ExposedPorts"`
}

// Option is a functional option for configuring the client.
type Option func(*Client)

// WithHost sets the host for the client.
func WithHost(host string) Option {
	return func(c *Client) {
		c.host = host
	}
}

// WithHost sets the basePath for the client.
func WithPath(path string) Option {
	return func(c *Client) {
		c.basePath = path + "v" + DefaultAPIVersion
	}
}

// NewClient creates a new Docker/Podman API client.
func NewClient(socketPath string, opts ...Option) (*Client, error) {
	c := &Client{
		basePath: "/v" + DefaultAPIVersion,
		host:     "localhost",
		scheme:   "http",
	}
	for _, opt := range opts {
		opt(c)
	}
	if strings.HasPrefix(socketPath, "unix://") {
		socketPath = strings.TrimPrefix(socketPath, "unix://")
		c.httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.DialTimeout("unix", socketPath, 30*time.Second)
				},
			},
			Timeout: 0,
		}
	} else if strings.HasPrefix(socketPath, "http://") || strings.HasPrefix(socketPath, "https://") {
		u, err := url.Parse(socketPath)
		if err != nil {
			return nil, fmt.Errorf("invalid host URL: %w", err)
		}
		c.scheme = u.Scheme
		c.host = u.Host
		c.httpClient = &http.Client{
			Timeout: 0,
		}
	} else {
		c.httpClient = &http.Client{
			Transport: &http.Transport{
				DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
					return net.DialTimeout("unix", socketPath, 30*time.Second)
				},
			},
			Timeout: 0,
		}
	}
	return c, nil
}

func (c *Client) buildURL(endpoint string, query url.Values) string {
	u := url.URL{
		Scheme:   c.scheme,
		Host:     c.host,
		Path:     c.basePath + endpoint,
		RawQuery: query.Encode(),
	}
	return u.String()
}

// doRequest performs an HTTP request and returns the response body.
// It handles common error cases like request creation, execution, and non-OK status codes.
// The caller is responsible for closing the returned ReadCloser.
func (c *Client) doRequest(ctx context.Context, method, endpoint string, query url.Values) (io.ReadCloser, error) {
	url := c.buildURL(endpoint, query)
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("executing request: %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
	return resp.Body, nil
}

// doQuery performs an HTTP request and decodes the JSON response into the provided result.
func (c *Client) doQuery(ctx context.Context, method, endpoint string, query url.Values, result interface{}) error {
	body, err := c.doRequest(ctx, method, endpoint, query)
	if err != nil {
		return err
	}
	defer body.Close()
	if result != nil {
		if err := json.NewDecoder(body).Decode(result); err != nil {
			return fmt.Errorf("decoding json response: %w", err)
		}
	}
	return nil
}

// ListContainers returns a list of running containers.
// Filters can be used to narrow down the results (e.g., by label).
func (c *Client) ListContainers(ctx context.Context, filters map[string][]string) ([]Container, error) {
	query := url.Values{}
	if len(filters) > 0 {
		filtersJSON, err := json.Marshal(filters)
		if err != nil {
			return nil, fmt.Errorf("encoding json filters: %w", err)
		}
		query.Set("filters", string(filtersJSON))
	}
	var containers []Container
	if err := c.doQuery(ctx, "GET", "/containers/json", query, &containers); err != nil {
		return nil, err
	}
	return containers, nil
}

// InspectContainer returns detailed information about a container.
func (c *Client) InspectContainer(ctx context.Context, containerID string) (*ContainerInspect, error) {
	var inspect ContainerInspect
	if err := c.doQuery(ctx, "GET", "/containers/"+containerID+"/json", nil, &inspect); err != nil {
		return nil, err
	}
	return &inspect, nil
}

// Events streams container events.
// Filters can be used to narrow down the events (e.g., by type, event, container).
func (c *Client) Events(ctx context.Context, filters map[string][]string) (<-chan Event, <-chan error) {
	eventCh := make(chan Event)
	errCh := make(chan error, 1)
	go func() {
		defer close(eventCh)
		defer close(errCh)
		query := url.Values{}
		if len(filters) > 0 {
			filtersJSON, err := json.Marshal(filters)
			if err != nil {
				errCh <- fmt.Errorf("encoding json filters: %w", err)
				return
			}
			query.Set("filters", string(filtersJSON))
		}
		body, err := c.doRequest(ctx, "GET", "/events", query)
		if err != nil {
			errCh <- err
			return
		}
		defer body.Close()
		scanner := bufio.NewScanner(body)
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 1024*1024)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				return
			default:
				line := scanner.Bytes()
				if len(line) == 0 {
					continue
				}
				var event Event
				if err := json.Unmarshal(line, &event); err != nil {
					continue
				}
				select {
				case eventCh <- event:
				case <-ctx.Done():
					return
				}
			}
		}
		if err := scanner.Err(); err != nil && ctx.Err() == nil {
			errCh <- fmt.Errorf("reading events: %w", err)
		}
	}()
	return eventCh, errCh
}

// Ping checks if the Docker/Podman daemon is accessible.
func (c *Client) Ping(ctx context.Context) error {
	return c.doQuery(ctx, "GET", "/_ping", nil, nil)
}
