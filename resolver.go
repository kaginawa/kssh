package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type report struct {
	ID            string `json:"id"`                  // MAC address
	CustomID      string `json:"custom_id,omitempty"` // User specified ID
	SSHServerHost string `json:"ssh_server_host"`     // Connected SSH server host
	SSHRemotePort int    `json:"ssh_remote_port"`     // Connected SSH remote port
	LocalIPv4     string `json:"ip4_local"`           // Local IPv4 address
	Hostname      string `json:"hostname"`            // OS Hostname
	ServerTime    int64  `json:"server_time"`         // Server-side consumed UTC time
}

type sshServer struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Key      string `json:"key"`
	Password string `json:"password"`
}

type client struct {
	server string
	apiKey string
}

func newClient(server, apiKey string) *client {
	c := client{
		server: server,
		apiKey: apiKey,
	}
	if !strings.HasPrefix(c.server, "http") {
		c.server = "https://" + c.server
	}
	return &c
}

func (c *client) findByID(mac string) (*report, error) {
	url := c.server + "/nodes/" + mac
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}
	req.Header.Set("Authorization", "token "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to communicate Kaginawa Server: %w", err)
	}
	defer safeClose(resp.Body, "connection")
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server respond HTTP %s", resp.Status)
	}
	var report report
	if err := json.NewDecoder(resp.Body).Decode(&report); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &report, nil
}

func (c *client) findByCustomID(cid string) ([]report, error) {
	url := c.server + "/nodes?custom-id=" + cid
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}
	req.Header.Set("Authorization", "token "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to communicate Kaginawa Server: %w", err)
	}
	defer safeClose(resp.Body, "connection")
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server respond HTTP %s", resp.Status)
	}
	var reports []report
	if err := json.NewDecoder(resp.Body).Decode(&reports); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return reports, nil
}

func (c *client) sshServer(host string) (*sshServer, error) {
	url := c.server + "/servers/" + host
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %w", err)
	}
	req.Header.Set("Authorization", "token "+c.apiKey)
	req.Header.Set("Accept", "application/json")
	resp, err := new(http.Client).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to communicate Kaginawa Server: %w", err)
	}
	defer safeClose(resp.Body, "connection")
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server respond HTTP %s", resp.Status)
	}
	var server sshServer
	if err := json.NewDecoder(resp.Body).Decode(&server); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &server, nil
}
