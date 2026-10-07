package requestrestmodel

import (
	"encoding/json"

	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
)

type UpdatePostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (r *UpdatePostRequest) FillFromBytes(body []byte) error {
	return json.Unmarshal(body, r)
}

func (r *UpdatePostRequest) ValidateRequest() error {
	if r.Title == "" {
		return restservermodel.ErrMissingRequiredField("title")
	}

	if r.Body == "" {
		return restservermodel.ErrMissingRequiredField("body")
	}

	return nil
}
