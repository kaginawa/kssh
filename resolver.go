package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

func (c *client) retrieve(url, method string, body io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, url, body)
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
	raw, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}
	return raw, nil
}

func (c *client) findByID(mac string) (*report, error) {
	body, err := c.retrieve(c.server+"/nodes/"+mac, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find %s: %w", mac, err)
	}
	var report report
	if err := json.Unmarshal(body, &report); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &report, nil
}

func (c *client) findByCustomID(cid string) ([]report, error) {
	body, err := c.retrieve(c.server+"/nodes?custom-id="+cid, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find %s: %w", cid, err)
	}
	var reports []report
	if err := json.Unmarshal(body, &reports); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return reports, nil
}

func (c *client) sshServer(host string) (*sshServer, error) {
	body, err := c.retrieve(c.server+"/servers/"+host, http.MethodGet, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to find %s: %w", host, err)
	}
	var server sshServer
	if err := json.Unmarshal(body, &server); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}
	return &server, nil
}
