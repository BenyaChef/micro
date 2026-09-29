package initinfra

import "micro/infrastructure/errors"

var ErrServiceNameIsRequired = errors.NewError("SYS", "InitInfra: ServiceName is required")

const defaultLogLevel = "info"
