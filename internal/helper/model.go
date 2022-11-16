package helper

import (
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

func ToUserResponse(user domain.User) web.UserResponse {
	return web.UserResponse{
		Id:        user.Id,
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Bio:       user.Bio,
		Email:     user.Email,
	}
}

func UserSetDefaultValue(userRequest *domain.User, userData *domain.User) {
	if userRequest.Username == "" {
		userRequest.Username = userData.Username
	}
	if userRequest.FirstName == "" {
		userRequest.FirstName = userData.FirstName
	}
	if userRequest.LastName == "" {
		userRequest.LastName = userData.LastName
	}
	if userRequest.Bio == "" {
		userRequest.Bio = userData.Bio
	}
	if userRequest.Email == "" {
		userRequest.Email = userData.Email
	}
}
