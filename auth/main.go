package main

import (
	"fmt"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/BenyaChef/micro/auth/initservice"
	"github.com/BenyaChef/micro/infrastructure/initinfra"
)

const serviceName = "auth"

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

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return container.RESTServer().Run(groupCtx)
	})

	group.Go(func() error {
		return container.GRPCServer().Run(groupCtx)
	})

	return group.Wait()
}
