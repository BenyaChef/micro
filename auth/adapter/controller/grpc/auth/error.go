package authgrpccontroller

import "github.com/BenyaChef/micro/infrastructure/errors"

var ErrTokenUseCaseIsRequired = errors.NewError("SYS", "AuthGRPCController: TokenUseCase is required")
