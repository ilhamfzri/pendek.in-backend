package router

import (
	"github.com/ilhamfzri/pendek.in/helper"
	"github.com/ilhamfzri/pendek.in/internal/controller"
	"github.com/ilhamfzri/pendek.in/internal/middleware"
)

func AddSocialMediaRoute(server *Server, socialMediaLinkController controller.SocialMediaLinkController, jwt helper.IJwt) {
	jwtMiddleware := middleware.NewJwtMiddleware(jwt.GetSigningKey())
	socialMediaRouteAuth := server.Router.Group("/v1/link/social-media")
	socialMediaRouteAuth.Use(jwtMiddleware)
	{
		socialMediaRouteAuth.GET("/types", socialMediaLinkController.GetAllTypes)
		socialMediaRouteAuth.POST("/", socialMediaLinkController.CreateLink)
		socialMediaRouteAuth.PUT("/:type_id", socialMediaLinkController.UpdateLink)
		socialMediaRouteAuth.GET("/", socialMediaLinkController.GetAllLink)
		socialMediaRouteAuth.GET("/analytic", socialMediaLinkController.GetLinkAnalytic)
		socialMediaRouteAuth.GET("/analytic/summary", socialMediaLinkController.GetSummaryLinkAnalytic)
	}
	socialMediaRoute := server.Router.Group("")
	socialMediaRoute.GET("/:username/:social-media", socialMediaLinkController.RedirectLink)

}
