package helper

import (
	"encoding/json"
	"net/http"
	"strings"
)

func GetAuthToken(request *http.Request) string {
	bearerToken := request.Header.Get("Authorization")
	splitToken := strings.Split(bearerToken, "Bearer ")
	jwtToken := splitToken[1]
	return jwtToken
}

func ReadFromRequestBody(request *http.Request, result interface{}) {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(&result)
	PanicIfError(err)
}

func WriteToResponse(writer http.ResponseWriter, status_code int, response interface{}) {
	writer.Header().Add("Content-Type", "application/json")
	writer.WriteHeader(status_code)
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(response)
	PanicIfError(err)
}
