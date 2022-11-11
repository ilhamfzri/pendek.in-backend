package main

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/ilhamfzri/pendek.in/app"
	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/ilhamfzri/pendek.in/internal/exception"
	"github.com/ilhamfzri/pendek.in/internal/helper"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"github.com/ilhamfzri/pendek.in/internal/service"
	_ "github.com/lib/pq"
)

func main() {
	db := app.NewDB()
	validator := validator.New()

	userRepository := repository.NewUserRepository()
	userService := service.NewUserService(userRepository, db, validator)
	userController := controller.NewUserController(userService)
	userRouter := app.NewUserRouter(userController)
	userRouter.PanicHandler = exception.ErrorHandler

	server := http.Server{
		Addr:    "localhost:3000",
		Handler: userRouter,
	}

	err := server.ListenAndServe()
	helper.PanicIfError(err)
}
