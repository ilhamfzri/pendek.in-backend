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
