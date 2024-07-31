package main

import (
	"fmt"
	"os"

	"github.com/lazyIoad/salo/internal/mnode"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name:  "saloserver",
		Usage: "salo managed node server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "socket",
				Value: "/tmp/salo/server.sock",
				Usage: "server socket path",
			},
		},
		Action: func(cCtx *cli.Context) {
			p := cCtx.String("socket")
			s := mnode.NewApiServer(p)
			s.Start()
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("failed to start worker server: %w", err))
		os.Exit(1)
	}
}
