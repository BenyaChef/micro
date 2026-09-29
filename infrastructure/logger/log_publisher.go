package logger

import (
	"context"
	"log/slog"

	apperrors "micro/infrastructure/errors"
)

type LogPublisher struct {
	log *slog.Logger
}

func (p *LogPublisher) LogError(ctx context.Context, errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}

		p.log.ErrorContext(ctx, err.Error(),
			slog.String("error_code", apperrors.CodeOf(err).String()),
		)
	}
}

func (p *LogPublisher) LogWarn(ctx context.Context, msg string, attrs ...any) {
	p.log.WarnContext(ctx, msg, attrs...)
}

func (p *LogPublisher) LogInfo(ctx context.Context, msg string, attrs ...any) {
	p.log.InfoContext(ctx, msg, attrs...)
}

func (p *LogPublisher) LogDebug(ctx context.Context, msg string, attrs ...any) {
	p.log.DebugContext(ctx, msg, attrs...)
}
