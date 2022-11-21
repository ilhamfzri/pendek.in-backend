package helper

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type JwtMock struct {
	mock.Mock
}

func (jwtClient *JwtMock) GetClaims(jwtToken string) JwtUserClaims {
	arguments := jwtClient.Mock.Called(jwtToken)
	return arguments.Get(0).(JwtUserClaims)
}

func (jwtClient *JwtMock) NewToken(id string, username string, email string) (string, time.Time, error) {
	arguments := jwtClient.Mock.Called(id, username, email)
	return arguments.String(0), arguments.Get(1).(time.Time), arguments.Error(2)
}

func (jwtClient *JwtMock) GetSigningKey() string {
	arguments := jwtClient.Mock.Called()
	return arguments.String(0)
}
