package commonautheventmodel

type UserEventAction string

const (
	UserEventActionUnknown    UserEventAction = ""
	UserEventActionRegistered UserEventAction = "REGISTERED"
)

func UserEventActionFrom(rawAction string) (UserEventAction, error) {
	action := UserEventAction(rawAction)

	switch action {
	case UserEventActionRegistered:
		return action, nil
	case UserEventActionUnknown:
		return UserEventActionUnknown, ErrActionIsEmpty
	default:
		return UserEventActionUnknown, ErrActionUnsupported(rawAction)
	}
}

func (a UserEventAction) String() string {
	return string(a)
}

func (a UserEventAction) IsZero() bool {
	return a == UserEventActionUnknown
}
