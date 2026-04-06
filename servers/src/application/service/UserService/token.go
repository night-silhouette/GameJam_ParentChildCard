package UserService

import (
	"context"
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

func (u *User_service_impl) Release_token(userID int, ctx context.Context) (string, global.ResponseStatusCode) {
	Active := u.repo.UpdateActiveInRedisByUserId(userID, ctx)
	expire_time := time.Now().Add(time.Hour * 24 * global.TokenExpiredTime)

	// 修改点：传入 ctx 和 u.repo.Get_db()
	e, err := u.repo.Get_by_id(ctx, u.repo.Get_db(), userID)
	if err != global.ResponseSuccess {
		return "", global.ResponseDataNotFound
	}

	claims := &Claims{
		UserId:   userID,
		Is_admin: e.Is_admin,
		Active:   Active,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire_time),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pcc_card_server",
			Subject:   "user_token",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(Key))
	return tokenString, global.ResponseSuccess
}

type Claims struct {
	UserId   int  `json:"user_id"`
	Is_admin bool `json:"is_admin"`
	jwt.RegisteredClaims
	Active int
}

func (u *User_service_impl) Is_valid_token(tokenString string, ctx context.Context) (int, bool, global.ResponseStatusCode) {
	var Flag global.ResponseStatusCode
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
	RealActive := u.repo.CheckActiveInRedisByUserId(id, ctx)
	is_admin := claims.Is_admin
	active := claims.Active
	if RealActive == active {
		Flag = global.ResponseSuccess
	} else {
		Flag = global.ResponseTokenHasUpdate
	}

	return id, is_admin, Flag
}
