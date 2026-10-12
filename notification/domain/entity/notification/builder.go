package notificationentity

import (
	"errors"

	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

type Builder struct {
	id      NotificationID
	eventID EventID
	userID  commonuserentity.UserID
	kind    Kind
	payload Payload
	sentAt  timeprimitive.Timestamp
}

func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) ID(id NotificationID) *Builder {
	b.id = id

	return b
}

func (b *Builder) EventID(eventID EventID) *Builder {
	b.eventID = eventID

	return b
}

func (b *Builder) UserID(userID commonuserentity.UserID) *Builder {
	b.userID = userID

	return b
}

func (b *Builder) Kind(kind Kind) *Builder {
	b.kind = kind

	return b
}

func (b *Builder) Payload(payload Payload) *Builder {
	b.payload = payload

	return b
}

func (b *Builder) SentAt(sentAt timeprimitive.Timestamp) *Builder {
	b.sentAt = sentAt

	return b
}

func (b *Builder) Build() (*Notification, error) {
	if err := b.checkRequiredFields(); err != nil {
		return nil, err
	}

	if b.sentAt.IsZero() {
		b.sentAt = timeprimitive.Now()
	}

	return &Notification{
		id:      b.id,
		eventID: b.eventID,
		userID:  b.userID,
		kind:    b.kind,
		payload: b.payload,
		sentAt:  b.sentAt,
	}, nil
}

func (b *Builder) checkRequiredFields() error {
	var errs []error

	if b.id.IsZero() {
		errs = append(errs, ErrNotificationIDIsRequired)
	}

	if b.eventID.IsZero() {
		errs = append(errs, ErrEventIDIsRequired)
	}

	if b.userID.IsZero() {
		errs = append(errs, ErrUserIDIsRequired)
	}

	if b.kind.IsZero() {
		errs = append(errs, ErrKindIsRequired)
	}

	if b.payload.IsZero() {
		errs = append(errs, ErrPayloadIsRequired)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
