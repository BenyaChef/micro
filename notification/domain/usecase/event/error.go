package eventusecase

import (
	"context"

	"github.com/BenyaChef/micro/infrastructure/errors"
	loggerinterface "github.com/BenyaChef/micro/infrastructure/logger/interface"
)

var (
	ErrHandlerProviderIsRequired = errors.NewError("SYS", "UserEventUseCase: HandlerProvider is required")
	ErrLogPublisherIsRequired    = errors.NewError("SYS", "UserEventUseCase: LogPublisher is required")
)

var ErrEventIsRequired = errors.NewErrorWithLevel(
	"0f96d603-001", "UserEventUseCase: event is required", errors.LevelInfo,
)

type errorProcessor struct {
	logPublisher loggerinterface.LogPublisher
}

func newErrorProcessor(logPublisher loggerinterface.LogPublisher) *errorProcessor {
	return &errorProcessor{logPublisher: logPublisher}
}

func (p *errorProcessor) LogAndReturn(ctx context.Context, err error) error {
	p.logPublisher.LogError(ctx, err)

	return err
}
