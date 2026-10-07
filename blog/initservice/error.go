package initservice

import "github.com/BenyaChef/micro/infrastructure/errors"

var ErrInfraContainerIsRequired = errors.NewError("SYS", "DependencyContainer: InfraContainer is required")
