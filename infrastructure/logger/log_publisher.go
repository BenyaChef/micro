package logger

import (
	"context"
	"log/slog"

	apperrors "github.com/BenyaChef/micro/infrastructure/errors"
)

type LogPublisher struct {
	log *slog.Logger
}

func (p *LogPublisher) LogError(ctx context.Context, errs ...error) {
	for _, err := range errs {
		if err == nil {
			continue
		}

		level := apperrors.LevelOf(err)
		attrs := []any{
			slog.String("error_code", apperrors.CodeOf(err).String()),
			slog.String("error_level", level.String()),
		}

		switch level {
		case apperrors.LevelInfo:
			p.log.InfoContext(ctx, err.Error(), attrs...)
		case apperrors.LevelWarn:
			p.log.WarnContext(ctx, err.Error(), attrs...)
		case apperrors.LevelError, apperrors.LevelCritical:
			p.log.ErrorContext(ctx, err.Error(), attrs...)
		default:
			p.log.ErrorContext(ctx, err.Error(), attrs...)
		}
	}
}

func (p *LogPublisher) LogErrorMsg(ctx context.Context, msg string, attrs ...any) {
	p.log.ErrorContext(ctx, msg, attrs...)
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
