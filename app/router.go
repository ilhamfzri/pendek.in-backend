package app

import (
	"fmt"
	"net/http"

	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/julienschmidt/httprouter"
)

func NewUserRouter(userController controller.UserController) *httprouter.Router {
	router := httprouter.New()
	router.GET("/v1/api/users/test", func(writer http.ResponseWriter, request *http.Request, params httprouter.Params) {
		fmt.Fprintf(writer, "pendek.in API")
	})
	router.POST("/v1/users/register", userController.Register)
	router.POST("/v1/users/login", userController.Login)
	router.GET("/v1/users/verify", userController.Verify)
	router.POST("/v1/users/change-password", userController.ChangePassword)
	router.PUT("/v1/users/update-info", userController.UpdateInformation)
	return router
}
