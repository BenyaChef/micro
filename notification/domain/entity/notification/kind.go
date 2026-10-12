package notificationentity

type Kind string

const (
	KindUnknown Kind = ""
	KindWelcome Kind = "welcome"
)

func KindFrom(rawKind string) (Kind, error) {
	kind := Kind(rawKind)

	switch kind {
	case KindWelcome:
		return kind, nil
	case KindUnknown:
		return KindUnknown, ErrKindIsRequired
	default:
		return KindUnknown, ErrKindUnsupported(rawKind)
	}
}

func (k Kind) String() string {
	return string(k)
}

func (k Kind) IsZero() bool {
	return k == KindUnknown
}
