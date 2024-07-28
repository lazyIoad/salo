package pipeline

import (
	"context"
	"fmt"

	"github.com/lazyIoad/salo/host"
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

func New(name string) *Pipeline {
	return &Pipeline{
		name: name,
	}
}

func (p *Pipeline) AddTask(name string, t taskif.Tasker) *Pipeline {
	p.tasks = append(p.tasks, task{name: name, impl: t})
	return p
}

func (p *Pipeline) Run(ctx context.Context, hosts []host.Host) error {
	fmt.Printf("PIPELINE [%s] starting...", p.name)
	for i, t := range p.tasks {
		fmt.Printf("TASK %d/%d [%s] starting...", i, len(p.tasks), t.name)
		t.impl.Run(ctx)
	}

	return nil
}
