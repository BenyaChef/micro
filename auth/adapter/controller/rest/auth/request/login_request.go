package requestrestmodel

import (
	"encoding/json"

	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (r *LoginRequest) FillFromBytes(body []byte) error {
	return json.Unmarshal(body, r)
}

func (r *LoginRequest) ValidateRequest() error {
	if r.Email == "" {
		return restservermodel.ErrMissingRequiredField("email")
	}

	if r.Password == "" {
		return restservermodel.ErrMissingRequiredField("password")
	}

	return nil
}
