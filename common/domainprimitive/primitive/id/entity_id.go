package idprimitive

import "github.com/google/uuid"

type EntityID string

func NewEntityID() EntityID {
	return EntityID(uuid.New().String())
}

func EntityIDFrom(rawID string) (EntityID, error) {
	if rawID == "" {
		return "", ErrEntityIDIsEmpty
	}

	parsedUUID, err := uuid.Parse(rawID)
	if err != nil {
		return "", ErrEntityIDInvalidFormat
	}

	return EntityID(parsedUUID.String()), nil
}

func (id EntityID) String() string {
	return string(id)
}
