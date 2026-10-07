package postentity

import idprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/id"

type PostID struct {
	value idprimitive.EntityID
}

func NewPostID() PostID {
	return PostID{value: idprimitive.NewEntityID()}
}

func PostIDFrom(rawID string) (PostID, error) {
	entityID, err := idprimitive.EntityIDFrom(rawID)
	if err != nil {
		return PostID{}, ErrPostIDInvalidFormat
	}

	return PostID{value: entityID}, nil
}

func (id PostID) String() string {
	return id.value.String()
}

func (id PostID) IsZero() bool {
	return id.value == ""
}
