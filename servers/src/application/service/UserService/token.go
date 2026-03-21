package UserService

import (
	"errors"
	"fmt"
	"pcc_card/global"
	"pcc_card/infra/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var Key string

func (u *User_service_impl) Init_key() {
	Key = config.Read_secret_key()
}

func (u *User_service_impl) Release_token(userID int) (string, global.ResponseStatusCode) {
	expire_time := time.Now().Add(time.Hour * 24 * global.TokenExpiredTime)
	e, err := u.repo.Get_by_id(userID)
	if err != global.ResponseSuccess {
		return "", global.ResponseDataNotFound
	}
	claims := &Claims{
		UserId:   userID,
		Is_admin: e.Is_admin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire_time),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pcc_card_server",
			Subject:   "user_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(Key))
	fmt.Println(tokenString)
	return tokenString, global.ResponseSuccess
}

type Claims struct {
	UserId   int  `json:"user_id"`
	Is_admin bool `json:"is_admin"`
	jwt.RegisteredClaims
}

func (u *User_service_impl) Is_valid_token(tokenString string) (int, bool, global.ResponseStatusCode) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		keyBytes := []byte(Key)
		return keyBytes, nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return 0, false, global.ResponseTokenExpired
		}
		if errors.Is(err, jwt.ErrSignatureInvalid) {
			return 0, false, global.ResponseInvalidToken
		}
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return 0, false, global.ResponseIncorrectTokenFormat
		}
	}
	id := claims.UserId
	is_admin := claims.Is_admin
	return id, is_admin, global.ResponseSuccess

}
