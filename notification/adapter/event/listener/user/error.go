package usereventlistener

import "github.com/BenyaChef/micro/infrastructure/errors"

var (
	ErrUserEventUseCaseIsRequired = errors.NewError("SYS", "UserEventListener: UserEventUseCase is required")
	ErrLogPublisherIsRequired     = errors.NewError("SYS", "UserEventListener: LogPublisher is required")
)

var ErrMessageIsRequired = errors.NewErrorWithLevel(
	"062b5dd8-001", "UserEventListener: message is required", errors.LevelInfo,
)
