package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const defaultConfigFileName = "kssh.conf"

type config struct {
	server string
	apiKey string
}

func loadConfig(path string) (config, error) {
	var c config
	if path == defaultConfigFileName {
		if defaultDir, err := os.UserConfigDir(); err == nil {
			path = filepath.Join(defaultDir, path)
		}
	}
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return c, fmt.Errorf("failed to load %s: %w", path, err)
	}
	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "AdminKey ") {
			c.apiKey = strings.TrimSpace(strings.Replace(text, "AdminKey ", "", 1))
			continue
		}
		if strings.HasPrefix(text, "Server ") {
			c.server = strings.TrimSpace(strings.Replace(text, "Server ", "", 1))
			continue
		}
	}
	return c, nil
}

func (c config) save(path string) error {
	if path == defaultConfigFileName {
		if defaultDir, err := os.UserConfigDir(); err == nil {
			if _, err := os.Stat(defaultDir); err != nil {
				if err := os.MkdirAll(defaultDir, 600); err != nil {
					return fmt.Errorf("failed to create directory: %w", err)
				}
			}
			path = filepath.Join(defaultDir, path)
		}
	}
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", path, err)
	}
	defer safeClose(f, path)
	content := fmt.Sprintf("Server %s\nAdminKey %s\n", c.server, c.apiKey)
	if _, err := f.WriteString(content); err != nil {
		return fmt.Errorf("failed to write %s: %w", path, err)
	}
	return nil
}
