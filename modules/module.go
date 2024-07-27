package modules

import "google.golang.org/grpc"

type Defaulter interface {
	Default() Service
}

type Service interface {
	RegisterServer(s grpc.ServiceRegistrar)
}
