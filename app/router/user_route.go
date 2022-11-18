package router

import (
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/ilhamfzri/pendek.in/internal/middleware"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"github.com/ilhamfzri/pendek.in/internal/service"
	"gorm.io/gorm"
)

func AddUsersRoute(server *Server, DB *gorm.DB, logger *logger.Logger, jwt *helper.Jwt) {

	userRepository := repository.NewUserRepository(logger)
	userService := service.NewUserService(userRepository, DB, logger, jwt)
	userController := controller.NewUserController(userService, logger)

	userRouteNotAuth := server.Router.Group("/v1/users")
	{
		userRouteNotAuth.POST("/sign-up", userController.Register)
		userRouteNotAuth.POST("/login", userController.Login)
		userRouteNotAuth.POST("/email-verification", userController.EmailVerification)
	}

	jwtMiddleware := middleware.NewJwtMiddleware(jwt.SigningKey)

	userRouteAuth := server.Router.Group("/v1/users")
	userRouteAuth.Use(jwtMiddleware)
	{
		userRouteAuth.POST("/change-password", userController.ChangePassword)
		userRouteAuth.PUT("/", userController.Update)
	}
}
