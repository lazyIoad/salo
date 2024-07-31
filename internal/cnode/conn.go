// Inspired heavily by github.com/johnsiilver/serveonssh
package cnode

import (
	"context"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

func init() {
	resolver.SetDefaultScheme("passthrough")
}

type SshProxiedGrpcConn struct {
	proxy *sshProxy
	*grpc.ClientConn
}

func NewSshProxiedGrpcConn(address string, port int, sshConfig *ssh.ClientConfig, socket string) (*SshProxiedGrpcConn, error) {
	proxy, err := newSshProxy(address, port, sshConfig, socket)
	if err != nil {
		return nil, fmt.Errorf("failed to create SSH proxy: %w", err)
	}

	opts := []grpc.DialOption{
		// We can be insecure here because this should all be encrypted by the
		// SSH tunnel anyway.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return proxy.dialer(ctx)
		}),
	}

	conn, err := grpc.NewClient(socket, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	return &SshProxiedGrpcConn{
		proxy,
		conn,
	}, nil
}

func (s *SshProxiedGrpcConn) Close() error {
	err := s.ClientConn.Close()
	if err != nil {
		return fmt.Errorf("failed to close gRPC connection: %w", err)
	}

	err = s.proxy.sshClient.Close()
	if err != nil {
		return fmt.Errorf("failed to close SSH connection: %w", err)
	}

	return nil
}

type dialer func(context.Context) (net.Conn, error)

type sshProxy struct {
	sshClient *ssh.Client
	dialer    dialer
}

func newSshProxy(address string, port int, sshConfig *ssh.ClientConfig, socket string) (*sshProxy, error) {
	addr := fmt.Sprintf("%s:%d", address, port)
	client, err := ssh.Dial("tcp", addr, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH host: %w", err)
	}

	var dial dialer = func(ctx context.Context) (net.Conn, error) {
		return client.DialContext(ctx, "unix", socket)
	}

	return &sshProxy{
		sshClient: client,
		dialer:    dial,
	}, nil
}
