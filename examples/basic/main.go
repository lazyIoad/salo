package main

import (
	"context"
	"log"
	"os"

	"github.com/lazyIoad/salo"
	"github.com/lazyIoad/salo/task/builtin/ping"
)

func main() {
	pw := os.Getenv("SSH_PASSWORD")
	cfg, err := salo.InsecureHostConfig(pw)
	fatalif(err)
	hosts := salo.NewHostsFromSlice(cfg, "localhost")

	p := salo.NewPipeline("ping").
		AddTask("Ping", ping.Default())

	ctx := context.Background()
	err = p.Execute(ctx, hosts)
	fatalif(err)
}

func fatalif(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
