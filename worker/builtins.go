package worker

import (
	"github.com/lazyIoad/salo/modules"
	"github.com/lazyIoad/salo/modules/ping"
	"google.golang.org/grpc"
)

var builtins = []modules.Defaulter{
	&ping.PingService{},
}

func registerBuiltins(s grpc.ServiceRegistrar) {
	for _, b := range builtins {
		b.Default().RegisterServer(s)
	}
}
