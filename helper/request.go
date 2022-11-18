package helper

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractTokenFromRequestHeader(c *gin.Context) string {
	bearerToken := c.Request.Header["Authorization"]
	splitToken := strings.Split(bearerToken[0], "Bearer ")
	jwtToken := splitToken[1]
	return jwtToken
}
