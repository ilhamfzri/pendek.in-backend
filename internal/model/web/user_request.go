package web

type UserRegisterRequest struct {
	Username string `json:"username" binding:"required,min=6,max=25,alphanum"`
	Email    string `json:"email" binding:"required,email,min=1,max=50"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email,min=1,max=50"`
	Password string `json:"password" binding:"required,min=6,max=16"`
}

type UserEmailVerificationRequest struct {
	Email            string `json:"email" binding:"required,email,min=1,max=50"`
	VerificationCode string `json:"verification_code" binding:"required,min=1,max=6"`
}

type UserChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required,min=6,max=16"`
	NewPassword     string `json:"new_password" binding:"required,min=6,max=16"`
}

type UserUpdateRequest struct {
	FullName string `validate:"max=16" json:"full_name,omitempty"`
	Bio      string `validate:"max=255" json:"bio,omitempty"`
}
