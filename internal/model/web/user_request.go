package web

type UserRegisterRequest struct {
	Username string `validate:"required,min=6,max=25" json:"username"`
	Email    string `validate:"required,email,min=1,max=50" json:"email"`
	Password string `validate:"required,min=6,max=16" json:"password"`
}

type UserLoginRequest struct {
	Email    string `validate:"required,email,min=1,max=50" json:"email"`
	Password string `validate:"required,min=6,max=16" json:"password"`
}

type UserVerifyRequest struct {
	Email string `validate:"required,email,min=1,max=50"`
	Code  string `validate:"required,min=1"`
}

type UserChangePasswordRequest struct {
	CurrentPassword string `validate:"required,min=6,max=16" json:"current_password"`
	NewPassword     string `validate:"required,min=6,max=16" json:"new_password"`
}

type UserUpdateInfoRequest struct {
	FirstName string `validate:"max=16" json:"first_name,omitempty"`
	LastName  string `validate:"max=16" json:"last_name,omitempty"`
	Bio       string `validate:"max=255" json:"bio,omitempty"`
}
