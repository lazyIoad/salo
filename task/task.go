package task

import "context"

type Tasker interface {
	Run(context.Context) error
	Plan() string
}
