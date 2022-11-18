package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ilhamfzri/pendek.in/app/logger"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

var InternalServerErrorResponse = web.WebResponseFailed{
	Status:  "failed",
	Message: "internal server error",
}

func NewRecoveryHandler(logger *logger.Logger) gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, InternalServerErrorResponse)
	})
}
