package responserestmodel

import boundarymodel "github.com/BenyaChef/micro/auth/boundary/model"

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func FromLoginResult(result *boundarymodel.LoginResult) *LoginResponse {
	return &LoginResponse{
		AccessToken: result.AccessToken,
		ExpiresIn:   result.ExpiresIn,
	}
}
