package restserverservice

import "github.com/BenyaChef/micro/infrastructure/errors"

var (
	ErrLogPublisherIsRequired         = errors.NewError("SYS", "ResponseService: LogPublisher is required")
	ErrErrorResponseServiceIsRequired = errors.NewError("SYS", "ResponseService: ErrorResponseService is required")
	ErrErrorResolverIsRequired        = errors.NewError("SYS", "ErrorResponseService: ErrorResolver is required")
	ErrErrorLogPublisherIsRequired    = errors.NewError("SYS", "ErrorResponseService: LogPublisher is required")
)
