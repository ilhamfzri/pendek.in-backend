package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"github.com/ilhamfzri/pendek.in/internal/helper"
	"github.com/ilhamfzri/pendek.in/internal/model/web"
)

var urlsNeedsJwt = []string{
	"/v1/users/change-password",
	"/v1/users/update-info"}

func urlNeedsJwt(urlPath string) bool {
	for _, url := range urlsNeedsJwt {
		if url == urlPath {
			return true
		}
	}
	return false
}

type JwtMiddleware struct {
	Handler    http.Handler
	SigningKey string
}

func NewJwtMiddleware(handler http.Handler, signingKey string) *JwtMiddleware {
	return &JwtMiddleware{
		Handler:    handler,
		SigningKey: signingKey,
	}
}

func (middleware *JwtMiddleware) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	urlPath := request.URL.Path
	if urlNeedsJwt(urlPath) {
		bearerToken := request.Header.Get("Authorization")
		if bearerToken == "" {
			webResponse := web.WebResponseFailed{
				Status:  "failed",
				Message: "authorization key not found",
			}
			helper.WriteToResponse(writer, http.StatusBadRequest, webResponse)
			return
		}

		splitToken := strings.Split(bearerToken, "Bearer ")
		if len(splitToken) != 2 {
			webResponse := web.WebResponseFailed{
				Status:  "failed",
				Message: "invalid bearer token format",
			}
			helper.WriteToResponse(writer, http.StatusBadRequest, webResponse)
			return
		}

		jwtToken := splitToken[1]
		token, errJwtParse := jwt.Parse(jwtToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(middleware.SigningKey), nil
		})

		if errJwtParse == jwt.ErrSignatureInvalid || !token.Valid {
			webResponse := web.WebResponseFailed{
				Status:  "failed",
				Message: "invalid or expired authentication key",
			}
			helper.WriteToResponse(writer, http.StatusUnauthorized, webResponse)
			return
		}
		helper.PanicIfError(errJwtParse)
	}
	middleware.Handler.ServeHTTP(writer, request)
}
