package mnode

import (
	"errors"
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/lazyIoad/salo/internal/modules"
	"google.golang.org/grpc"
)

type ApiServer struct {
	socketAddr string
	modules    []modules.Service
}

func NewApiServer(socketAddr string) *ApiServer {
	return &ApiServer{
		socketAddr: socketAddr,
		modules:    make([]modules.Service, 0),
	}
}

func (s *ApiServer) WithModule(m modules.Service) {
	s.modules = append(s.modules, m)
}

func (s *ApiServer) Start() {
	if err := os.MkdirAll(filepath.Dir(s.socketAddr), 0700); err != nil {
		log.Fatalf("failed to create socket path: %s", err)
	}

	if err := os.Remove(s.socketAddr); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			log.Fatalf("failed to remove existing socket: %s", err)
		}
	}

	l, err := net.Listen("unix", s.socketAddr)
	if err != nil {
		log.Fatalf("failed to listen on socket: %v", err)
	}

	g := grpc.NewServer()

	// First register all the builtins
	registerBuiltins(g)

	// Then register any custom modules
	for _, m := range s.modules {
		m.RegisterServer(g)
	}

	if err := g.Serve(l); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
