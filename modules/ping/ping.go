package ping

import (
	"context"

	"github.com/lazyIoad/salo/modules"
	pb "github.com/lazyIoad/salo/modules/ping/proto"
	"google.golang.org/grpc"
)

type PingService struct {
	pb.UnimplementedPingerServer
}

func (p *PingService) Default() modules.Service {
	return &PingService{}
}

func (p *PingService) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Data: req.GetData(),
	}, nil
}

func (p *PingService) RegisterServer(s grpc.ServiceRegistrar) {
	pb.RegisterPingerServer(s, p)
}
