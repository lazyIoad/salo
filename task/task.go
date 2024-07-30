package task

import (
	"context"

	"google.golang.org/grpc"
)

type Tasker interface {
	Run(context.Context, grpc.ClientConnInterface) error
}
