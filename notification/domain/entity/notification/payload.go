package notificationentity

import (
	"encoding/json"
)

type Payload []byte

func PayloadFrom(raw []byte) (Payload, error) {
	if len(raw) == 0 || !json.Valid(raw) {
		return nil, ErrPayloadIsEmpty
	}

	return Payload(raw), nil
}

func (p Payload) Bytes() []byte {
	return []byte(p)
}

func (p Payload) String() string {
	return string(p)
}

func (p Payload) IsZero() bool {
	return len(p) == 0
}
