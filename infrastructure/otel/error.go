package otel

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodeResourceFailed = errors.ErrorCode("32b339bb-010")
	ErrCodeExporterFailed = errors.ErrorCode("32b339bb-011")
	ErrCodeMonitorFailed  = errors.ErrorCode("32b339bb-012")
	ErrCodeListenFailed   = errors.ErrorCode("32b339bb-040")
)

var (
	ErrServiceNameIsRequired  = errors.NewError("SYS", "OTelProvider: ServiceName is required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "OTelProvider: LogPublisher is required")
	ErrPortIsRequired         = errors.NewError("SYS", "MetricsServer: Port is required")
)

func ErrResourceFailed(cause error) error {
	return errors.NewErrorWithLevelAndCause(
		ErrCodeResourceFailed, "Unable to build OTel resource", errors.LevelCritical, cause,
	)
}

func ErrExporterFailed(cause error) error {
	return errors.NewErrorWithLevelAndCause(
		ErrCodeExporterFailed, "Unable to create Prometheus exporter", errors.LevelCritical, cause,
	)
}

func ErrMonitorFailed(name string, cause error) error {
	msg := fmt.Sprintf("Unable to create monitor: name=%s", name)

	return errors.NewErrorWithLevelAndCause(ErrCodeMonitorFailed, msg, errors.LevelCritical, cause)
}

func ErrListenFailed(addr string, cause error) error {
	msg := fmt.Sprintf("Metrics server failed to listen: addr=%s", addr)

	return errors.NewErrorWithLevelAndCause(ErrCodeListenFailed, msg, errors.LevelCritical, cause)
}
