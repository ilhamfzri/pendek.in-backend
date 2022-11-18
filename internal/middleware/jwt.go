package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/ilhamfzri/pendek.in/helper"
)

var ErrAuthNotFound = errors.New("authorization key not found")
var ErrInvalidBearerToken = errors.New("invalid bearer token format")
var ErrInvalidOrExpiredToken = errors.New("invalid or expired authentication key")

func NewJwtMiddleware(signingKey string) gin.HandlerFunc {
	return func(c *gin.Context) {

		// Checking if the request has an authorization header. If not, it will return an error.
		bearerToken := c.Request.Header["Authorization"]
		if len(bearerToken) != 1 {
			c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(ErrAuthNotFound))
			return
		}

		// Checking if the token is in the correct format.
		splitToken := strings.Split(bearerToken[0], "Bearer ")
		if len(splitToken) != 2 {
			c.AbortWithStatusJSON(http.StatusBadRequest, helper.ToWebResponseFailed(ErrInvalidBearerToken))
			return
		}

		// Checking if the token is valid or not.
		jwtToken := splitToken[1]
		token, errJwtParse := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(signingKey), nil
		})

		if errJwtParse != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.ToWebResponseFailed(ErrInvalidOrExpiredToken))
			return
		}

		c.Next()
	}
}
