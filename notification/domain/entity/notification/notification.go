package notificationentity

import (
	commonuserentity "github.com/BenyaChef/micro/common/auth/domain/entity/user"
	timeprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/time"
)

type Notification struct {
	id      NotificationID
	eventID EventID
	userID  commonuserentity.UserID
	kind    Kind
	payload Payload
	sentAt  timeprimitive.Timestamp
}

func (n *Notification) ID() NotificationID {
	return n.id
}

func (n *Notification) EventID() EventID {
	return n.eventID
}

func (n *Notification) UserID() commonuserentity.UserID {
	return n.userID
}

func (n *Notification) Kind() Kind {
	return n.kind
}

func (n *Notification) Payload() Payload {
	return n.payload
}

func (n *Notification) SentAt() timeprimitive.Timestamp {
	return n.sentAt
}
