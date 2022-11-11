package controller

import (
	"context"
	"net/http"

	"github.com/ilhamfzri/pendek.in/internal/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
	"github.com/ilhamfzri/pendek.in/internal/service"
	"github.com/julienschmidt/httprouter"
)

type UserControllerImpl struct {
	UserService service.UserService
}

func NewUserController(userService service.UserService) UserController {
	return &UserControllerImpl{
		UserService: userService,
	}
}

func (controller *UserControllerImpl) Register(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userRegisterRequest := web.UserRegisterRequest{}
	helper.ReadFromRequestBody(request, &userRegisterRequest)

	ctx := context.Background()
	userResponse, errService := controller.UserService.Register(ctx, userRegisterRequest)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		helper.WriteToResponse(writer, http.StatusBadRequest, webResponse)

	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "success create a new account",
			Data:    userResponse,
		}
		helper.WriteToResponse(writer, http.StatusCreated, webResponse)
	}

}

func (controller *UserControllerImpl) Login(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userLoginRequest := web.UserLoginRequest{}
	helper.ReadFromRequestBody(request, &userLoginRequest)

	ctx := context.Background()
	tokenResponse, errService := controller.UserService.Login(ctx, userLoginRequest)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		helper.WriteToResponse(writer, http.StatusBadRequest, webResponse)

	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "login success",
			Data:    tokenResponse,
		}
		helper.WriteToResponse(writer, http.StatusCreated, webResponse)
	}
}

func (controller *UserControllerImpl) Verify(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
	userVerifyRequest := web.UserVerifyRequest{}

	userVerifyRequest.Email = request.URL.Query().Get("email")
	userVerifyRequest.Code = request.URL.Query().Get("code")

	ctx := context.Background()
	errService := controller.UserService.Verify(ctx, userVerifyRequest)

	if errService != nil {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: errService.Error(),
		}
		helper.WriteToResponse(writer, http.StatusBadRequest, webResponse)

	} else {
		webResponse := web.WebResponseSuccess{
			Status:  "success",
			Message: "email verification success",
		}
		helper.WriteToResponse(writer, http.StatusOK, webResponse)
	}
}
