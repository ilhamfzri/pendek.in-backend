package helper

import (
	"github.com/ilhamfzri/pendek.in/internal/model/domain"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

func UserDomainToResponse(user *domain.User) web.UserResponse {
	return web.UserResponse{
		ID:       user.ID,
		Username: user.Username,
		FullName: user.FullName,
		Bio:      user.Bio,
		Email:    user.Email,
	}
}
