package main

import (
	"fmt"
	"os"

	envregistrymodel "github.com/BenyaChef/micro/infrastructure/envregistry/model"
	"github.com/BenyaChef/micro/infrastructure/initinfra"
	"github.com/BenyaChef/micro/infrastructure/restserver"
)

const (
	serviceName = "notification"
	defaultPort = "8083"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	container, err := initinfra.Bootstrap(serviceName)
	if err != nil {
		return err
	}

	env := container.Env()

	server, err := restserver.NewBuilder().
		LogPublisher(container.LogPublisher()).
		ServiceName(serviceName).
		Port(env.StringDefault(envregistrymodel.KeyRestPort, defaultPort)).
		Build()
	if err != nil {
		return err
	}
	if err := env.Err(); err != nil {
		return err
	}

	ctx, stop := initinfra.SignalContext()
	defer stop()

	return server.Run(ctx)
}
