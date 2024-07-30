package salo

import (
	"context"
	"fmt"
	"sync"
)

func executePipeline(ctx context.Context, p *Pipeline, hosts []*Host) error {
	var wg sync.WaitGroup

	for _, h := range hosts {
		wg.Add(1)

		go func() {
			defer wg.Done()
			hexec, err := newHostExecutor(h)
			if err != nil {
				panic(err)
			}

			defer hexec.close()
			err = hexec.runPipeline(ctx, p)
			if err != nil {
				panic(err)
			}
		}()
	}

	wg.Wait()
	return nil
}

type hostExecutor struct {
	host *Host
	conn *sshProxiedGrpcConn
}

func newHostExecutor(h *Host) (*hostExecutor, error) {
	conn, err := newSshProxiedGrpcConn(h)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize connection: %w", err)
	}

	return &hostExecutor{
		host: h,
		conn: conn,
	}, nil
}

func (h *hostExecutor) runPipeline(ctx context.Context, p *Pipeline) error {
	// TODO: cannot use normal logging
	fmt.Printf("[%s] PIPELINE (%s) starting...", h.host.Address, p.name)

	for i, t := range p.tasks {
		fmt.Printf("[%s] TASK [%d/%d: %s] starting...", h.host.Address, i, len(p.tasks), t.name)
		err := t.impl.Run(ctx, h.conn)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *hostExecutor) close() error {
	return h.conn.Close()
}
