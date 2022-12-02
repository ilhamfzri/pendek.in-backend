package router

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/ilhamfzri/pendek.in/internal/handler"
	"github.com/ilhamfzri/pendek.in/internal/middleware"
)

func AddUsersRoute(server *Server, userController controller.UserController, jwt helper.IJwt) {
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
		userRouteAuth.GET("/", userController.GetCurrentProfile)

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

	userRoutePublic := server.Router.Group("")
	userRoutePublic.GET("/:username", userController.Profile)

}
