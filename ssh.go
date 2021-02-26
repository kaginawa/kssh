package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/kaginawa/kaginawa-sdk-go"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

const eofRetries = 3

func connect(tunnel *kaginawa.SSHServer, user, defaultPassword string, port int) *ssh.Session {
	tunnelConfig, err := createSSHConfig(tunnel.User, tunnel.Key, tunnel.Password)
	if err != nil {
		fatalf("failed to create SSH config: %v", err)
	}
	var session *ssh.Session
	password := defaultPassword
	eofCount := 0
	for i := 0; ; i++ {
		// Connect to SSH tunneling server
		tConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", tunnel.Host, tunnel.Port), tunnelConfig)
		if err != nil {
			fatalf("failed to connect SSH tunneling server: %v", err)
		}

		// Connect to target
		conn, err := tConn.Dial("tcp", fmt.Sprintf("localhost:%d", port))
		if err != nil {
			fatalf("failed to connect: %v", err)
		}
		if i == 0 && len(defaultPassword) == 0 {
			password = passwordPrompt(user)
		}
		sshConfig, err := createSSHConfig(user, "", password)
		if err != nil {
			fatalf("failed to create SSH config: %v", err)
		}
		c, nc, req, err := ssh.NewClientConn(conn, fmt.Sprintf("localhost:%d", port), sshConfig)
		if err != nil {
			if strings.HasSuffix(err.Error(), "EOF") {
				safeClose(conn, "tcp connection")
				safeClose(tConn, "tunnel connection")
				eofCount++
				if eofCount >= eofRetries {
					fatalf("EOF occurred %d times", eofCount)
				}
				continue
			} else if strings.Contains(err.Error(), "unable to authenticate") {
				fmt.Printf("%v\n", err.Error())
				password = passwordPrompt(user)
				continue
			}
			fatalf("failed to create tunneling connection: %v", err)
		}
		client := ssh.NewClient(c, nc, req)

		// Open session
		session, err = client.NewSession()
		if err != nil {
			if strings.HasSuffix(err.Error(), "EOF") {
				safeClose(client, "ssh client")
				safeClose(tConn, "tunnel connection")
				continue // retry
			}
			fatalf("failed to create client session: %v", err)
		}
		break
	}
	return session
}

func createSSHConfig(user, key, password string) (*ssh.ClientConfig, error) {
	config := ssh.ClientConfig{
		User:            user,
		Auth:            make([]ssh.AuthMethod, 0),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	if len(key) > 0 {
		parsed, err := ssh.ParsePrivateKey([]byte(key))
		if err != nil {
			return nil, err
		}
		config.Auth = append(config.Auth, ssh.PublicKeys(parsed))
	}
	if len(password) > 0 {
		config.Auth = append(config.Auth, ssh.Password(password))
	}
	return &config, nil
}

func passwordPrompt(user string) string {
	fmt.Printf("Password for %s: ", user)
	password, err := terminal.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fatalf("failed to read password: %v", err)
	}
	fmt.Println()
	return string(password)
}

func openTerminal(session *ssh.Session) {
	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		fmt.Println(err)
	}
	defer func() {
		if err := terminal.Restore(fd, state); err != nil {
			handleError(0, err)
		}
	}()
	w, h, err := terminal.GetSize(fd)
	if err != nil {
		fatalf("%v", err)
	}
	modes := ssh.TerminalModes{
		ssh.ECHO:          1,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	session.Stdout = os.Stdout
	session.Stderr = os.Stderr
	session.Stdin = os.Stdin
	if err := session.RequestPty("xterm", h, w, modes); err != nil {
		fatalf("failed to request terminal: %v", err)
	}
	if err := session.Shell(); err != nil {
		fatalf("failed to start shell: %s", err)
	}
	if err := session.Wait(); err != nil {
		fatalf("session brake: %v", err)
	}
}

func execCommand(s *ssh.Session, cmd string) (string, error) {
	result, err := s.Output(cmd)
	return string(result), err
}
