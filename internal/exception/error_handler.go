package exception

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/ilhamfzri/pendek.in/internal/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

func ErrorHandler(writer http.ResponseWriter, request *http.Request, err interface{}) {
	if validationErrors(writer, request, err) {
		return
	}
	internalServerError(writer, request, err)
}

func validationErrors(writer http.ResponseWriter, request *http.Request, err interface{}) bool {
	exception, ok := err.(validator.ValidationErrors)
	if ok {
		webResponse := web.WebResponseFailed{
			Status:  "failed",
			Message: exception.Error(),
		}
		helper.WriteToResponse(writer, http.StatusBadRequest, webResponse)
		return true
	} else {
		return false
	}
}

func internalServerError(writer http.ResponseWriter, request *http.Request, err interface{}) {
	fmt.Println(err)
	webResponse := web.WebResponseFailed{
		Status:  "failed",
		Message: "internal server error",
	}
	helper.WriteToResponse(writer, http.StatusInternalServerError, webResponse)
}
