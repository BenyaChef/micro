package commonuserentity

import idprimitive "github.com/BenyaChef/micro/common/domainprimitive/primitive/id"

type UserID struct {
	value idprimitive.EntityID
}

func NewUserID() UserID {
	return UserID{value: idprimitive.NewEntityID()}
}

func UserIDFrom(rawID string) (UserID, error) {
	entityID, err := idprimitive.EntityIDFrom(rawID)
	if err != nil {
		return UserID{}, ErrUserIDInvalidFormat
	}

	return UserID{value: entityID}, nil
}

func (id UserID) String() string {
	return id.value.String()
}

func (id UserID) IsZero() bool {
	return id.value == ""
}
