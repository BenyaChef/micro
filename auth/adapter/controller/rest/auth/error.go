package authrestcontroller

import "github.com/BenyaChef/micro/infrastructure/errors"

var (
	ErrBaseControllerIsRequired = errors.NewError("SYS", "AuthController: BaseController is required")
	ErrUserUseCaseIsRequired    = errors.NewError("SYS", "AuthController: UserUseCase is required")
	ErrTokenUseCaseIsRequired   = errors.NewError("SYS", "AuthController: TokenUseCase is required")
)

const ErrCodeAuthorizationRequired = errors.ErrorCode("8b8c0e30-001")

var ErrAuthorizationRequired = errors.NewErrorWithLevel(
	ErrCodeAuthorizationRequired, "Authorization header with Bearer token is required", errors.LevelInfo,
)
