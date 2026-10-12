package notificationentity

import (
	"fmt"

	"github.com/BenyaChef/micro/infrastructure/errors"
)

const ErrCodeKindUnsupported = errors.ErrorCode("bd49bce2-004")

var (
	ErrNotificationIDInvalidFormat = errors.NewErrorWithLevel(
		"bd49bce2-001", "Notification id has invalid format", errors.LevelInfo,
	)
	ErrEventIDInvalidFormat = errors.NewErrorWithLevel(
		"bd49bce2-002", "Event id has invalid format", errors.LevelInfo,
	)
	ErrPayloadIsEmpty = errors.NewErrorWithLevel(
		"bd49bce2-003", "Notification payload is empty", errors.LevelInfo,
	)
)

var (
	ErrNotificationIDIsRequired = errors.NewError("bd49bce2-010", "Notification: NotificationID is required")
	ErrEventIDIsRequired        = errors.NewError("bd49bce2-011", "Notification: EventID is required")
	ErrUserIDIsRequired         = errors.NewError("bd49bce2-012", "Notification: UserID is required")
	ErrKindIsRequired           = errors.NewError("bd49bce2-013", "Notification: Kind is required")
	ErrPayloadIsRequired        = errors.NewError("bd49bce2-014", "Notification: Payload is required")
)

func ErrKindUnsupported(rawKind string) error {
	msg := fmt.Sprintf("Notification kind is not supported: kind=%s", rawKind)

	return errors.NewErrorWithLevel(ErrCodeKindUnsupported, msg, errors.LevelWarn)
}
