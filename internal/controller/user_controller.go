package controller

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type UserController interface {
	Register(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	Verify(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	ChangePassword(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
	UpdateInformation(writer http.ResponseWriter, request *http.Request, params httprouter.Params)
}
