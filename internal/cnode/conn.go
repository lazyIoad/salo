// Inspired heavily by github.com/johnsiilver/serveonssh
package cnode

import (
	"context"
	"fmt"
	"net"

	"github.com/lazyIoad/salo/host"
	"golang.org/x/crypto/ssh"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type SshProxiedGrpcConn struct {
	proxy *sshProxy
	*grpc.ClientConn
}

func NewSshProxiedGrpcConn(h *host.Host) (*SshProxiedGrpcConn, error) {
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

	return &SshProxiedGrpcConn{
		proxy,
		conn,
	}, nil
}

func (s *SshProxiedGrpcConn) Close() error {
	return s.proxy.sshClient.Close()
}

type dialer func(context.Context) (net.Conn, error)

type sshProxy struct {
	sshClient *ssh.Client
	dialer    dialer
}

func newSshProxy(h *host.Host) (*sshProxy, error) {
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
