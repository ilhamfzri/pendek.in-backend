package helper

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func NewJwt(signingKey string, expiredTimeDay int, issuer string) *Jwt {
	return &Jwt{
		SigningKey:     signingKey,
		ExpiredTimeDay: expiredTimeDay,
		Issuer:         issuer,
	}
}

type Jwt struct {
	SigningKey     string
	ExpiredTimeDay int
	Issuer         string
}

type JwtUserClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	jwt.StandardClaims
}

func (jwtClient *Jwt) GetClaims(jwtToken string) JwtUserClaims {
	jwtClaims := JwtUserClaims{}
	_, err := jwt.ParseWithClaims(jwtToken, &jwtClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtClient.SigningKey), nil
	})
	PanicIfError(err)
	return jwtClaims
}

func (jwtClient *Jwt) NewToken(username string, email string) (string, time.Time) {
	expiredTime := time.Now().Add(time.Hour * 24 * time.Duration(jwtClient.ExpiredTimeDay))
	jwtClaims := JwtUserClaims{
		username,
		email,
		jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
			Issuer:    jwtClient.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	key, err := token.SignedString([]byte(jwtClient.SigningKey))
	PanicIfError(err)
	return key, expiredTime
}
