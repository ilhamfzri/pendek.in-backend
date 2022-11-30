package helper

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func GenerateCacheKeyByJwt(c *gin.Context) string {
	requestUri := c.Request.URL.RequestURI()
	jwt := ExtractTokenFromRequestHeader(c)
	fmt.Println(jwt + requestUri)
	return jwt + requestUri
}
