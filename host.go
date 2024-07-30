package salo

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

type Host struct {
	Address string
	Port    int
	Config  *HostConfig
}

type HostConfig struct {
	// Path to create a unix socket for the API. If left empty, will default to
	// /tmp/salo/server.sock
	SocketPath string
	// SSH config to control how the SSH tunnel to the host will be established.
	// If left empty, will default to using the SSH agent and connect to the
	// same user as the invoking one.
	SshConfig *ssh.ClientConfig
}

// Creates a slice of hosts for the given addresses. The same config
// will be used for all hosts.
func NewHostsFromSlice(config *HostConfig, addrs ...string) []*Host {
	var hosts []*Host

	for _, a := range addrs {
		host := &Host{
			Address: a,
			Port:    22,
			Config:  config,
		}

		hosts = append(hosts, host)
	}

	return hosts
}

func DefaultHostConfig() (*HostConfig, error) {
	sshConfig, err := defaultSshSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to build SSH settings: %w", err)
	}

	return buildHostConfig(sshConfig)
}

func InsecureHostConfig(password string) (*HostConfig, error) {
	sshConfig := insecureSshSettings(password)
	return buildHostConfig(sshConfig)
}

func buildHostConfig(sshConfig *ssh.ClientConfig) (*HostConfig, error) {
	user := os.Getenv("USER")
	if user == "" {
		return nil, ErrUserNotFound
	}

	sshConfig.User = user

	return &HostConfig{
		SocketPath: "/tmp/salo/server.sock",
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
