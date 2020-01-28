package main

import (
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

var (
	configFile = flag.String("c", defaultConfigFileName, "path to configuration file")
	apiKey     = flag.String("k", "", "admin API key for the Kaginawa Server")
	server     = flag.String("s", "", "hostname of the Kaginawa Server")
)

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(2)
	}
	target := flag.Arg(0)
	if *configFile == defaultConfigFileName {
		if defaultDir, err := os.UserConfigDir(); err == nil {
			*configFile = filepath.Join(defaultDir, *configFile)
		}
	}
	username := ""
	if strings.Count(target, "@") == 1 {
		split := strings.Split(target, "@")
		username = split[0]
		target = split[1]
	} else {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		username = u.Name
	}
	var config config

	// Order of parameter sources
	// 1. Flags (-k and -s)
	// 2. Configuration file (-c or default)
	// 3. Interactive prompt (save after collection)
	if len(*apiKey) > 0 && len(*server) > 0 {
		config.apiKey = *apiKey
		config.server = *server
	} else if _, err := os.Stat(*configFile); err == nil {
		c, err := loadConfig(*configFile)
		if err != nil {
			fmt.Printf("failed to load %s: %v\n", *configFile, err)
		}
		config = c
	} else {
		config = inputConfig(*configFile)
	}

	var report report
	// Resolve by MAC
	client := newClient(config.server, config.apiKey)
	if strings.Count(target, ":") == 5 {
		r, err := client.findByID(target)
		if err != nil {
			fatalf("%v", err)
		}
		if r == nil {
			fatalf("target not found: %s", target)
		}
		report = *r
	} else {
		// Resolve by custom ID
		reports, err := client.findByCustomID(target)
		if err != nil {
			fatalf("%v", err)
		}
		if len(reports) == 0 {
			fatalf("target not found: %s", target)
		}
		report = selectTarget(reports)
	}
	if report.SSHRemotePort == 0 {
		fatalf("ssh not connected.")
	}
	tunnel, err := client.sshServer(report.SSHServerHost)
	if err != nil {
		fatalf("failed to get ssh server information: %v", err)
	}
	if tunnel == nil {
		fatalf("unknown ssh server: %s", report.SSHServerHost)
	}
	connect(tunnel, username, report.SSHRemotePort)
}

func inputConfig(path string) config {
	fmt.Printf("Creating configuration file: %s\n", *configFile)
	var c config
	for {
		fmt.Print("Kaginawa Server (ex. xxx.yyy.com): ")
		if _, err := fmt.Scan(&c.server); err != nil {
			continue
		}
		c.server = strings.TrimSpace(c.server)
		if len(c.server) == 0 {
			continue
		}
		break
	}
	for {
		fmt.Print("Admin API key: ")
		if _, err := fmt.Scan(&c.apiKey); err != nil {
			continue
		}
		c.apiKey = strings.TrimSpace(c.apiKey)
		if len(c.apiKey) == 0 {
			continue
		}
		break
	}
	if err := c.save(path); err != nil {
		fatalf("failed to create %s: %v", path, err)
	}
	return c
}

func selectTarget(reports []report) report {
	if len(reports) == 1 {
		return reports[0]
	}
	fmt.Printf("Multiple choice:\n")
	for i, r := range reports {
		fmt.Printf("%d: %s %s %s", i+1, r.ID, r.LocalIPv4, r.Hostname)
	}
	var n int
	for {
		fmt.Print("number > ")
		if _, err := fmt.Scan(&n); err != nil {
			continue
		}
		n = n - 1
		if n < 0 || n >= len(reports) {
			fmt.Println("out of range")
			continue
		}
		return reports[n]
	}
}
