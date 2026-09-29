package logger

import "micro/infrastructure/errors"

var ErrServiceNameIsRequired = errors.NewError("SYS", "LogPublisher: ServiceName is required")
