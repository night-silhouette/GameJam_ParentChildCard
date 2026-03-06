package UserService

import (
	"errors"
	"pcc_card/global"
	"pcc_card/infra/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Key string

func (u *User_service_impl) Init_key() {
	Key = config.Read_secret_key()
}

func (u *User_service_impl) Release_token(userID int) string {
	expire_time := time.Now().Add(time.Hour * 24 * 30)

	claims := &Claims{
		UserId: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire_time),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pcc_card_server",
			Subject:   "user_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString(Key)
	return tokenString
}

type Claims struct {
	UserId int `json:"user_id"`
	jwt.RegisteredClaims
}

func (u *User_service_impl) Is_valid_token(tokenString string) (int, global.ResponseStatusCode) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		keyBytes := []byte(Key)
		return keyBytes, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, global.ResponseTokenExpired
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return 0, global.ResponseInvalidToken
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return 0, global.ResponseIncorrectTokenFormat
		}
	}
	id := claims.UserId
	return id, global.ResponseSuccess

}
