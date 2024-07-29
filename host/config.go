package host

import (
	"errors"
	"fmt"
	"net"
	"os"
	"path"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

var (
	ErrUserNotFound           = errors.New("could not determine user from $USER")
	ErrSshAgentSocketNotFound = errors.New("could not find SSH agent socket from $SSH_AUTH_SOCK - is it running?")
)

type Config struct {
	// Path to create a unix socket for the API. If left empty, will default to
	// /tmp/salo/saloserver.sock
	SocketPath string
	// SSH config to control how the SSH tunnel to the host will be established.
	// If left empty, will default to using the SSH agent and connect to the
	// same user as the invoking one.
	SshConfig *ssh.ClientConfig
}

func DefaultConfig() (*Config, error) {
	sshConfig, err := defaultSshSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to build SSH settings: %w", err)
	}

	return buildConfig(sshConfig)
}

func InsecureConfig(password string) (*Config, error) {
	sshConfig := insecureSshSettings(password)
	return buildConfig(sshConfig)
}

func buildConfig(sshConfig *ssh.ClientConfig) (*Config, error) {
	user := os.Getenv("USER")
	if user == "" {
		return nil, ErrUserNotFound
	}

	sshConfig.User = user

	return &Config{
		SocketPath: "/tmp/salo/saloserver.sock",
		SshConfig:  sshConfig,
	}, nil
}

func defaultSshSettings() (*ssh.ClientConfig, error) {
	sp := os.Getenv("SSH_AUTH_SOCK")
	if sp == "" {
		return nil, ErrSshAgentSocketNotFound
	}

	conn, err := net.Dial("unix", sp)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to ssh-agent: %w", err)
	}

	client := agent.NewClient(conn)
	auth := ssh.PublicKeysCallback(client.Signers)

	cb, err := getHostKeyCallback()
	if err != nil {
		return nil, fmt.Errorf("failed to parse known_hosts file: %w", err)
	}

	return &ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{auth},
		HostKeyCallback: cb,
	}, nil
}

func insecureSshSettings(password string) *ssh.ClientConfig {
	return &ssh.ClientConfig{
		Auth:            []ssh.AuthMethod{ssh.Password(password)}, // not great
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),              // bad
	}
}

func getHostKeyCallback() (ssh.HostKeyCallback, error) {
	knownHostsPath := path.Join(os.Getenv("HOME"), ".ssh", "known_hosts")
	// TODO: this requires the host to already be in known_hosts. We should
	// support adding a new host on the fly without erroring.
	return knownhosts.New(knownHostsPath)
}
