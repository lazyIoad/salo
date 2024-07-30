// Inspired heavily by github.com/johnsiilver/serveonssh
package salo

import (
	"context"
	"fmt"
	"net"

	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type sshProxiedGrpcConn struct {
	proxy *sshProxy
	*grpc.ClientConn
}

func newSshProxiedGrpcConn(h *Host) (*sshProxiedGrpcConn, error) {
	proxy, err := newSshProxy(h)
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

	conn, err := grpc.NewClient(h.Config.SocketPath, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC client: %w", err)
	}

	return &sshProxiedGrpcConn{
		proxy,
		conn,
	}, nil
}

func (s *sshProxiedGrpcConn) close() error {
	return s.proxy.sshClient.Close()
}

type dialer func(context.Context) (net.Conn, error)

type sshProxy struct {
	sshClient *ssh.Client
	dialer    dialer
}

func newSshProxy(h *Host) (*sshProxy, error) {
	addr := fmt.Sprintf("%s:%d", h.Address, h.Port)
	client, err := ssh.Dial("tcp", addr, h.Config.SshConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to dial SSH host: %w", err)
	}

	var dial dialer = func(ctx context.Context) (net.Conn, error) {
		return client.DialContext(ctx, "unix", h.Config.SocketPath)
	}

	return &sshProxy{
		sshClient: client,
		dialer:    dial,
	}, nil
}
