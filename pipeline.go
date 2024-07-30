package salo

import (
	"context"

	taskif "github.com/lazyIoad/salo/task"
)

type task struct {
	name string
	impl taskif.Tasker
}

type Pipeline struct {
	name  string
	tasks []task
}

func NewPipeline(name string) *Pipeline {
	return &Pipeline{
		name: name,
	}
}

func (p *Pipeline) AddTask(name string, t taskif.Tasker) *Pipeline {
	p.tasks = append(p.tasks, task{name: name, impl: t})
	return p
}

func (p *Pipeline) Execute(ctx context.Context, hosts []*Host) error {
	return executePipeline(ctx, p, hosts)
}
