package queueserver

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodeConnectFailed = errors.ErrorCode("021072fc-040")
	ErrCodeSendFailed    = errors.ErrorCode("021072fc-041")
)

var (
	ErrBrokersAreRequired     = errors.NewError("SYS", "QueuePublisher: Brokers are required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "QueuePublisher: LogPublisher is required")
)

var ErrMessageIsRequired = errors.NewErrorWithLevel(
	"021072fc-001", "QueuePublisher: message is required", errors.LevelInfo,
)

func ErrConnectFailed(brokers string, cause error) error {
	msg := fmt.Sprintf("Unable to connect to Kafka: brokers=%s", brokers)

	return errors.NewErrorWithLevelAndCause(ErrCodeConnectFailed, msg, errors.LevelCritical, cause)
}

func ErrSendFailed(topic string, cause error) error {
	msg := fmt.Sprintf("Unable to publish message: topic=%s", topic)

	return errors.NewErrorWithCause(ErrCodeSendFailed, msg, cause)
}
