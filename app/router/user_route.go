package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/ilhamfzri/pendek.in/internal/handler"
	"github.com/ilhamfzri/pendek.in/internal/middleware"
	"github.com/ilhamfzri/pendek.in/internal/repository"
	"github.com/ilhamfzri/pendek.in/internal/service"
	"gorm.io/gorm"
)

func AddUsersRoute(server *Server, DB *gorm.DB, logger *logger.Logger, jwt helper.IJwt) {

	userRepository := repository.NewUserRepository(logger)
	userService := service.NewUserService(userRepository, DB, logger, jwt)
	userController := controller.NewUserController(userService, logger)

	userRouteNotAuth := server.Router.Group("/v1/users")
	{
		userRouteNotAuth.POST("/sign-up", userController.Register)
		userRouteNotAuth.POST("/login", userController.Login)
		userRouteNotAuth.POST("/email-verification", userController.EmailVerification)
	}

	jwtMiddleware := middleware.NewJwtMiddleware(jwt.GetSigningKey())
	userRouteAuth := server.Router.Group("/v1/users")
	userRouteAuth.Use(jwtMiddleware)
	{
		userRouteAuth.POST("/change-picture", userController.ChangeProfilePicture)
		userRouteAuth.GET("/generate-token", userController.GenerateToken)
		userRouteAuth.POST("/change-password", userController.ChangePassword)
		userRouteAuth.PUT("/", userController.Update)

	}

	userResourcePath := os.Getenv("PROFILE_IMG_DIR")
	userResourceDirectory := http.Dir(userResourcePath)
	userRouteResources := server.Router.Group("/v1/resources/users")

	// add middleware for security reason to prevent user see all files inside filesystem
	userRouteResources.Use(func(c *gin.Context) {
		urlPath := c.Request.URL.Path
		if strings.HasSuffix(urlPath, "/") {
			c.AbortWithStatusJSON(http.StatusNotFound, handler.NoRouteResponse)
		}
		c.Next()
	})
	{
		userRouteResources.StaticFS("/pictures/", userResourceDirectory)
	}

}
