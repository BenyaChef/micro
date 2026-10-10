package commonautheventmodel

import (
	"encoding/json"
	"errors"
	"time"
)

type UserEventPayload struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type UserEvent struct {
	EventID    string           `json:"event_id"`
	Action     UserEventAction  `json:"action"`
	OccurredAt time.Time        `json:"occurred_at"`
	Data       UserEventPayload `json:"data"`
}

func NewUserRegistered(eventID, userID, email string, occurredAt time.Time) *UserEvent {
	return &UserEvent{
		EventID:    eventID,
		Action:     UserEventActionRegistered,
		OccurredAt: occurredAt.UTC(),
		Data:       UserEventPayload{UserID: userID, Email: email},
	}
}

func UserEventFromBytes(raw []byte) (*UserEvent, error) {
	event := &UserEvent{}
	if err := json.Unmarshal(raw, event); err != nil {
		return nil, ErrUnmarshalFailed(err)
	}

	if err := event.Validate(); err != nil {
		return nil, err
	}

	return event, nil
}

func (e *UserEvent) ToBytes() ([]byte, error) {
	return json.Marshal(e)
}

func (e *UserEvent) PartitionKey() string {
	return e.Data.UserID
}

func (e *UserEvent) Validate() error {
	var errs []error

	if e.EventID == "" {
		errs = append(errs, ErrEventIDIsEmpty)
	}

	if _, err := UserEventActionFrom(e.Action.String()); err != nil {
		errs = append(errs, err)
	}

	if e.Data.UserID == "" {
		errs = append(errs, ErrUserIDIsEmpty)
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
