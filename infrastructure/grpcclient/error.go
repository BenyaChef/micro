package grpcclient

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const ErrCodeDialFailed = errors.ErrorCode("59ef47c7-040")

var (
	ErrTargetIsRequired       = errors.NewError("SYS", "GRPCClient: Target is required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "GRPCClient: LogPublisher is required")
)

func ErrDialFailed(target string, cause error) error {
	msg := fmt.Sprintf("Unable to create gRPC client: target=%s", target)

	return errors.NewErrorWithLevelAndCause(ErrCodeDialFailed, msg, errors.LevelCritical, cause)
}
