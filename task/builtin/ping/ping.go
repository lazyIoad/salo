package ping

import (
	"context"
	"fmt"

	pb "github.com/lazyIoad/salo/internal/modules/ping/proto"
	"google.golang.org/grpc"
)

type Ping struct {
	data string
}

func Default() *Ping {
	return &Ping{
		data: "greetings",
	}
}

func (p *Ping) Run(ctx context.Context, conn grpc.ClientConnInterface) error {
	client := pb.NewPingerClient(conn)
	resp, err := client.Ping(ctx, &pb.PingRequest{
		Data: p.data,
	})

	if err != nil {
		return fmt.Errorf("received error from node: %w", err)
	}

	fmt.Println(resp.Data)
	return nil
}
