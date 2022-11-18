package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

var NoMethodResponse = web.WebResponseFailed{
	Status:  "failed",
	Message: "error no method",
}

func NewNoMethodHandler() gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		c.AbortWithStatusJSON(http.StatusInternalServerError, NoMethodResponse)
	})
}
