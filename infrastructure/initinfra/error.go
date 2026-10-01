package initinfra

import "github.com/BenyaChef/micro/infrastructure/errors"

var ErrServiceNameIsRequired = errors.NewError("SYS", "InitInfra: ServiceName is required")

const defaultLogLevel = "info"
