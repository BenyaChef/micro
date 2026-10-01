package restservercontroller

import "github.com/BenyaChef/micro/infrastructure/errors"

var (
	ErrResponseServiceIsRequired = errors.NewError("SYS", "BaseController: ResponseService is required")
	ErrLogPublisherIsRequired    = errors.NewError("SYS", "BaseController: LogPublisher is required")
)
