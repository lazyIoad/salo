package builtin

import "context"

type Ping struct {
	data string
}

func Default() *Ping {
	return &Ping{
		data: "pong",
	}
}

func (p *Ping) Run(context.Context) error { return nil }

func (p *Ping) Plan() string { return "" }
