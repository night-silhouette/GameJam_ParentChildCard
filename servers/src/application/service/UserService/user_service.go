package UserService

import (
	"pcc_card/application/entity"
	"pcc_card/application/service"
	"pcc_card/global"
	"pcc_card/infra/repo"

	"golang.org/x/crypto/bcrypt"
)

type User_service interface {
	service.Service
	Find_user_by_name(name string) (*entity.User, global.ResponseStatusCode)
	Find_user_by_id(id int) (*entity.User, global.ResponseStatusCode)
	Check_password(password, hash string) bool
	Release_token(userID int) string
	Is_valid_token(tokenString string) (int, global.ResponseStatusCode)
	Create_user(user *entity.User) global.ResponseStatusCode
	Delete_user(e *entity.User) global.ResponseStatusCode
}

type User_service_impl struct {
	repo repo.User_repo
}

func (u *User_service_impl) Set_repo(r repo.Repo) {
	u.repo = r.(repo.User_repo)
	u.Init_key()
}
func (u *User_service_impl) Find_user_by_name(name string) (*entity.User, global.ResponseStatusCode) {
	e, err := u.repo.Get_by_name(name)
	if err != global.ResponseSuccess {
		return &entity.User{}, err
	} else {
		return e, global.ResponseSuccess
	}
}
func (u *User_service_impl) Find_user_by_id(id int) (*entity.User, global.ResponseStatusCode) {
	e, err := u.repo.Get_by_id(id)
	if err != global.ResponseSuccess {
		return &entity.User{}, err
	} else {
		return e, global.ResponseSuccess
	}
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

func (u *User_service_impl) Create_user(user *entity.User) global.ResponseStatusCode {
	user.Password = u.hash_password(user.Password)
	return u.repo.Create(user)
}

func (u *User_service_impl) Delete_user(e *entity.User) global.ResponseStatusCode {
	err := u.repo.Delete(e)
	if err != global.ResponseSuccess {
		return err
	}
	return global.ResponseSuccess
}
