package grpcserver

import (
	"fmt"
	"time"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodeListenFailed    = errors.ErrorCode("6c09ac91-010")
	ErrCodeServeFailed     = errors.ErrorCode("6c09ac91-011")
	ErrCodeShutdownTimeout = errors.ErrorCode("6c09ac91-012")
)

var (
	ErrLogPublisherIsRequired = errors.NewError("SYS", "GRPCServer: LogPublisher is required")
	ErrPortIsRequired         = errors.NewError("SYS", "GRPCServer: Port is required")
)

func ErrListenFailed(addr string, cause error) error {
	msg := fmt.Sprintf("gRPC server failed to listen: addr=%s", addr)

	return errors.NewErrorWithLevelAndCause(ErrCodeListenFailed, msg, errors.LevelCritical, cause)
}

func ErrServeFailed(addr string, cause error) error {
	msg := fmt.Sprintf("gRPC server failed to serve: addr=%s", addr)

	return errors.NewErrorWithLevelAndCause(ErrCodeServeFailed, msg, errors.LevelCritical, cause)
}

func ErrShutdownTimeout(addr string, timeout time.Duration) error {
	msg := fmt.Sprintf("gRPC server did not stop in time: addr=%s, timeout=%s", addr, timeout)

	return errors.NewError(ErrCodeShutdownTimeout, msg)
}
