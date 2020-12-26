package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kaginawa/kaginawa-sdk-go"
)

var (
	ver        = "v0.0.0"
	configFile = flag.String("c", defaultConfigFileName, "path to configuration file")
	apiKey     = flag.String("k", "", "admin API key for the Kaginawa Server")
	server     = flag.String("s", "", "hostname of the Kaginawa Server")
	v          = flag.Bool("v", false, "print version")
)

func main() {
	flag.Parse()
	if *v {
		fmt.Printf("kssh %s, compiled by %s\n", ver, runtime.Version())
		os.Exit(0)
	}
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
		config = configPrompt(*configFile)
	}

	var username, defaultPassword string
	if strings.Count(target, "@") == 1 {
		split := strings.Split(target, "@")
		username = split[0]
		target = split[1]
	} else if len(config.defaultUser) > 0 {
		username = config.defaultUser
		defaultPassword = config.defaultPassword
	} else {
		u, err := user.Current()
		if err != nil {
			panic(err)
		}
		username = u.Name
	}

	start(config, target, username, defaultPassword)
}

func configPrompt(path string) config {
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

func selectTarget(reports []kaginawa.Report) kaginawa.Report {
	if len(reports) == 1 {
		return reports[0]
	}
	fmt.Printf("Multiple choices:\n")
	for i, r := range reports {
		if len(r.LocalIPv4) == 0 && len(r.LocalIPv6) > 0 {
			fmt.Printf("%d: %s %s@%s %s\n", i+1, r.ID, r.LocalIPv6, r.Adapter, r.Hostname)
		} else {
			fmt.Printf("%d: %s %s@%s %s\n", i+1, r.ID, r.LocalIPv4, r.Adapter, r.Hostname)
		}
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

func start(config config, target string, username, defaultPassword string) {
	endpoint := config.server
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		endpoint = "https://" + endpoint
	}
	var report kaginawa.Report
	client, err := kaginawa.NewClient(endpoint, config.apiKey)
	if err != nil {
		fatalf("failed to prepare API client: %v", err)
	}
	if strings.Count(target, ":") == 5 {
		r, err := client.FindNode(context.Background(), target)
		if err != nil {
			fatalf("%v", err)
		}
		if r == nil {
			fatalf("target not found: %s", target)
		}
		report = *r
	} else {
		reports, err := client.ListNodesByCustomID(context.Background(), target)
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
	tunnel, err := client.FindSSHServerByHostname(context.Background(), report.SSHServerHost)
	if err != nil {
		fatalf("failed to get ssh server information: %v", err)
	}
	if tunnel == nil {
		fatalf("unknown ssh server: %s", report.SSHServerHost)
	}
	connect(tunnel, username, defaultPassword, report.SSHRemotePort)
}
