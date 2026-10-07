package postrestcontroller

import "github.com/BenyaChef/micro/infrastructure/errors"

const ErrCodeInvalidQueryParam = errors.ErrorCode("42601feb-001")

var (
	ErrBaseControllerIsRequired = errors.NewError("SYS", "PostController: BaseController is required")
	ErrPostUseCaseIsRequired    = errors.NewError("SYS", "PostController: PostUseCase is required")
)

var ErrCallerIsRequired = errors.NewErrorWithLevel(
	"42601feb-002", "PostController: caller is not authenticated", errors.LevelWarn,
)

func ErrInvalidQueryParam(param, value string) error {
	return errors.NewErrorWithLevel(
		ErrCodeInvalidQueryParam, "Invalid query parameter: "+param+"="+value, errors.LevelInfo,
	)
}
