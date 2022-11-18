package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

var NoRouteResponse = web.WebResponseFailed{
	Status:  "failed",
	Message: "error no route",
}

func NewNoRouteHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusNotFound, NoRouteResponse)
	})
}
