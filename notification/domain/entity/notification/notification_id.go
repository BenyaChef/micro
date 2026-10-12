package notificationentity

import idprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/id"

type NotificationID struct {
	value idprimitive.EntityID
}

func NewNotificationID() NotificationID {
	return NotificationID{value: idprimitive.NewEntityID()}
}

func NotificationIDFrom(rawID string) (NotificationID, error) {
	entityID, err := idprimitive.EntityIDFrom(rawID)
	if err != nil {
		return NotificationID{}, ErrNotificationIDInvalidFormat
	}

	return NotificationID{value: entityID}, nil
}

func (id NotificationID) String() string {
	return id.value.String()
}

func (id NotificationID) IsZero() bool {
	return id.value == ""
}
