package web

type UserRegisterRequest struct {
	FirstName string `validate:"required,max=50,min=1" json:"first_name"`
	LastName  string `validate:"required,max=50,min=1" json:"last_name"`
	Username  string `validate:"required,max=50,min=1" json:"username"`
	Email     string `validate:"required,max=255,min=1" json:"email"`
	Password  string `validate:"required,max=255,min=1" json:"password"`
}
