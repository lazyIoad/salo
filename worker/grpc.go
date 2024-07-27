package worker

import (
	"fmt"
	"log"
	"net"

	"github.com/lazyIoad/salo/modules"
	"google.golang.org/grpc"
)

type ApiServer struct {
	port    int
	modules []modules.Service
}

func NewApiServer(port int) *ApiServer {
	return &ApiServer{
		port:    port,
		modules: make([]modules.Service, 0),
	}
}

func (s *ApiServer) WithModule(m modules.Service) {
	s.modules = append(s.modules, m)
}

func (s *ApiServer) Start() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	g := grpc.NewServer()

	// First register all the builtins
	registerBuiltins(g)

	// Then register any custom modules
	for _, m := range s.modules {
		m.RegisterServer(g)
	}

	if err := g.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
