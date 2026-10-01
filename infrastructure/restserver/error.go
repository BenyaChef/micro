package restserver

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

var (
	ErrLoggerIsRequired      = errors.NewError("SYS", "RESTServer: Logger is required")
	ErrServiceNameIsRequired = errors.NewError("SYS", "RESTServer: ServiceName is required")
	ErrPortIsRequired        = errors.NewError("SYS", "RESTServer: Port is required")
)

const (
	ErrCodeListenFailed   = errors.ErrorCode("60ac51cc-010")
	ErrCodeShutdownFailed = errors.ErrorCode("60ac51cc-011")
)

func ErrListenFailed(addr string, cause error) error {
	msg := fmt.Sprintf("REST server failed to listen: addr=%s", addr)

	return errors.NewErrorWithLevelAndCause(ErrCodeListenFailed, msg, errors.LevelCritical, cause)
}

func ErrShutdownFailed(addr string, cause error) error {
	msg := fmt.Sprintf("REST server failed to shutdown: addr=%s", addr)

	return errors.NewErrorWithCause(ErrCodeShutdownFailed, msg, cause)
}
