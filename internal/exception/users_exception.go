package exception

type UserExceptionError struct {
	Error string
}

func NewUserExceptionError(error string) UserExceptionError {
	return UserExceptionError{Error: error}
}
