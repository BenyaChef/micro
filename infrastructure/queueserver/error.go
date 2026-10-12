package queueserver

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const (
	ErrCodeConnectFailed     = errors.ErrorCode("021072fc-040")
	ErrCodeSendFailed        = errors.ErrorCode("021072fc-041")
	ErrCodeConsumerGroupInit = errors.ErrorCode("021072fc-042")
	ErrCodeConsumeFailed     = errors.ErrorCode("021072fc-043")
)

var (
	ErrBrokersAreRequired     = errors.NewError("SYS", "QueuePublisher: Brokers are required")
	ErrLogPublisherIsRequired = errors.NewError("SYS", "QueuePublisher: LogPublisher is required")
)

var (
	ErrConsumerBrokersAreRequired     = errors.NewError("SYS", "QueueConsumer: Brokers are required")
	ErrConsumerLogPublisherIsRequired = errors.NewError("SYS", "QueueConsumer: LogPublisher is required")
	ErrConsumerGroupIDIsRequired      = errors.NewError("SYS", "QueueConsumer: GroupID is required")
	ErrConsumerTopicsAreRequired      = errors.NewError("SYS", "QueueConsumer: Topics are required")
	ErrConsumerHandlerIsRequired      = errors.NewError("SYS", "QueueConsumer: Handler is required")
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

func ErrConsumerGroupInit(groupID string, cause error) error {
	msg := fmt.Sprintf("Unable to create consumer group: groupID=%s", groupID)

	return errors.NewErrorWithLevelAndCause(ErrCodeConsumerGroupInit, msg, errors.LevelCritical, cause)
}

func ErrConsumeFailed(topics string, cause error) error {
	msg := fmt.Sprintf("Consumer session failed: topics=%s", topics)

	return errors.NewErrorWithCause(ErrCodeConsumeFailed, msg, cause)
}
