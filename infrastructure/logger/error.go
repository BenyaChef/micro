package logger

import "github.com/BenyaChef/micro/infrastructure/errors"

var ErrServiceNameIsRequired = errors.NewError("SYS", "LogPublisher: ServiceName is required")
