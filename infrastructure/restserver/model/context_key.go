package restservermodel

import "context"

type contextKey string

const (
	requestIDContextKey contextKey = "request_id"
	callerIDContextKey  contextKey = "caller_id"
)

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, requestIDContextKey, requestID)
}

func RequestIDFrom(ctx context.Context) string {
	requestID, ok := ctx.Value(requestIDContextKey).(string)
	if !ok {
		return ""
	}

	return requestID
}

func WithCallerID(ctx context.Context, callerID string) context.Context {
	return context.WithValue(ctx, callerIDContextKey, callerID)
}

func CallerIDFrom(ctx context.Context) string {
	callerID, ok := ctx.Value(callerIDContextKey).(string)
	if !ok {
		return ""
	}

	return callerID
}
