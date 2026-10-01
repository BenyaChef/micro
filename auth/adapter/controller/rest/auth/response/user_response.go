package responserestmodel

import userentity "github.com/BenyaChef/micro/auth/domain/entity/user"

type UserResponse struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

func FromUser(user *userentity.User) *UserResponse {
	return &UserResponse{
		UserID: user.ID().String(),
		Email:  user.Email().String(),
	}
}
