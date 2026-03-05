package UserService

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"pcc_card/global"
	"pcc_card/infra/config"
	"time"
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

func (u *User_service_impl) Is_valid_token(tokenString string) (int, global.StatusCode) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		keyBytes := []byte(Key)
		return keyBytes, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, global.StatusTokenExpired
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return 0, global.StatusInvalidToken
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return 0, global.StatusIncorrectTokenFormat
		}
	}
	id := claims.UserId
	return id, global.StatusSuccess

}
