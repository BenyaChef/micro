package main

import (
	"fmt"
	"os"

	"github.com/BenyaChef/micro/blog/initservice"
	"github.com/BenyaChef/micro/infrastructure/initinfra"
)

const serviceName = "blog"

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	infra, err := initinfra.Bootstrap(serviceName)
	if err != nil {
		return err
	}

	ctx, stop := initinfra.SignalContext()
	defer stop()

	container, err := initservice.NewDependencyContainer(ctx, infra, serviceName)
	if err != nil {
		return err
	}

	defer container.Close(ctx)

	return container.RESTServer().Run(ctx)
}
