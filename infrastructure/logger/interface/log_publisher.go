package loggerinterface

import "context"

type LogPublisher interface {
	LogError(ctx context.Context, errs ...error)
	LogErrorMsg(ctx context.Context, msg string, attrs ...any)
	LogWarn(ctx context.Context, msg string, attrs ...any)
	LogInfo(ctx context.Context, msg string, attrs ...any)
	LogDebug(ctx context.Context, msg string, attrs ...any)
}
