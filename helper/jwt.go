package helper

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/ilhamfzri/pendek.in/config"
)

type IJwt interface {
	GetClaims(jwtToken string) JwtUserClaims
	NewToken(id string, username string, email string) (string, time.Time, error)
	GetSigningKey() string
}

type Jwt struct {
	SigningKey     string
	ExpiredTimeDay int
	Issuer         string
}

func NewJwt(cfg config.JwtConfig) IJwt {
	return &Jwt{
		SigningKey:     cfg.SigningKey,
		ExpiredTimeDay: cfg.ExpiredTimeDay,
		Issuer:         cfg.Issuer,
	}
}

type JwtUserClaims struct {
	Id       string `json:"id"`
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

func (jwtClient *Jwt) NewToken(id string, username string, email string) (string, time.Time, error) {
	expiredTime := time.Now().Add(time.Hour * 24 * time.Duration(jwtClient.ExpiredTimeDay))
	jwtClaims := JwtUserClaims{
		id,
		username,
		email,
		jwt.StandardClaims{
			ExpiresAt: expiredTime.Unix(),
			Issuer:    jwtClient.Issuer,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	key, err := token.SignedString([]byte(jwtClient.SigningKey))
	return key, expiredTime, err
}

func (jwtClient *Jwt) GetSigningKey() string {
	return jwtClient.SigningKey
}
