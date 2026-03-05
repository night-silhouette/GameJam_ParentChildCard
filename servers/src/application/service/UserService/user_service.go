package UserService

import (
	"database/sql"
	"errors"
	"fmt"
	"pcc_card/application/entity"
	"pcc_card/application/service"
	"pcc_card/global"
	"pcc_card/infra/repo"

	"golang.org/x/crypto/bcrypt"
)

type User_service interface {
	service.Service
	Find_user_by_name(name string) (*entity.User, global.StatusCode)
	Find_user_by_id(id int) (*entity.User, global.StatusCode)
	Check_password(password, hash string) bool
	Release_token(userID int) string
	Is_valid_token(tokenString string) (int, global.StatusCode)
}

type User_service_impl struct {
	repo repo.User_repo
}

func (u *User_service_impl) Set_repo(r repo.Repo) {
	u.repo = r.(repo.User_repo)
	u.Init_key()
}
func (u *User_service_impl) Find_user_by_name(name string) (*entity.User, global.StatusCode) {
	e, err := u.repo.Get_by_name(name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &entity.User{}, global.StatusDataNotFound
		} else {
			fmt.Println("sql err", err)
			return &entity.User{}, global.StatusInternalServersError
		}
	}
	return e, global.StatusSuccess
}
func (u *User_service_impl) Find_user_by_id(id int) (*entity.User, global.StatusCode) {
	e, err := u.repo.Get_by_id(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &entity.User{}, global.StatusDataNotFound
		} else {
			fmt.Println("sql err", err)
			return &entity.User{}, global.StatusInternalServersError
		}
	}
	return e, global.StatusSuccess
}
func (u *User_service_impl) hash_password(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
func (u *User_service_impl) Check_password(password, hash string) bool {
	if u.hash_password(password) == hash {
		return true
	} else {
		return false
	}
}
