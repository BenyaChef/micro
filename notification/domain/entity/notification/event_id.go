package notificationentity

import idprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/id"

type EventID struct {
	value idprimitive.EntityID
}

func EventIDFrom(rawID string) (EventID, error) {
	entityID, err := idprimitive.EntityIDFrom(rawID)
	if err != nil {
		return EventID{}, ErrEventIDInvalidFormat
	}

	return EventID{value: entityID}, nil
}

func (id EventID) String() string {
	return id.value.String()
}

func (id EventID) IsZero() bool {
	return id.value == ""
}
