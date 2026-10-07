package requestrestmodel

import (
	"encoding/json"

	restservermodel "github.com/BenyaChef/micro/infrastructure/restserver/model"
)

type CreatePostRequest struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (r *CreatePostRequest) FillFromBytes(body []byte) error {
	return json.Unmarshal(body, r)
}

func (r *CreatePostRequest) ValidateRequest() error {
	if r.Title == "" {
		return restservermodel.ErrMissingRequiredField("title")
	}

	if r.Body == "" {
		return restservermodel.ErrMissingRequiredField("body")
	}

	return nil
}
