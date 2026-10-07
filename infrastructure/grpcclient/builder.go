package grpcclient

import (
	"errors"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	grpcclientmodel "github.com/BenyaChef/micro/infrastructure/grpcclient/model"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

type Builder struct {
	target       string
	logPublisher loggerinterface.LogPublisher
	callTimeout  time.Duration
}

func NewBuilder() *Builder {
	return &Builder{callTimeout: grpcclientmodel.DefaultCallTimeout}
}

func (b *Builder) Target(target string) *Builder {
	b.target = target

	return b
}

func (b *Builder) LogPublisher(logPublisher loggerinterface.LogPublisher) *Builder {
	b.logPublisher = logPublisher

	return b
}

func (b *Builder) CallTimeout(timeout time.Duration) *Builder {
	b.callTimeout = timeout

	return b
}

func (b *Builder) Build() (*Client, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	connection, err := grpc.NewClient(
		b.target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(callTimeout(b.callTimeout)),
	)
	if err != nil {
		return nil, ErrDialFailed(b.target, err)
	}

	return &Client{connection: connection, logPublisher: b.logPublisher, target: b.target}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.target == "" {
		errs = append(errs, ErrTargetIsRequired)
	}

	if b.logPublisher == nil {
		errs = append(errs, ErrLogPublisherIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
