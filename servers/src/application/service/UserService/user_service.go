package UserService

import (
	"context"
	"fmt"
	"pcc_card/application/entity/User_entity"
	"pcc_card/application/service"
	"pcc_card/global"
	"pcc_card/infra/repo"
	"pcc_card/infra/repo/userrepo"

	"golang.org/x/crypto/bcrypt"
)

type User_service interface {
	service.Service
	Find_user_by_name(name string) (*User_entity.User, global.ResponseStatusCode)
	Find_user_by_id(id int) (*User_entity.User, global.ResponseStatusCode)
	Check_password(id int, password string) global.ResponseStatusCode
	Release_token(userID int, ctx context.Context) (string, global.ResponseStatusCode)
	Is_valid_token(tokenString string, ctx context.Context) (int, bool, global.ResponseStatusCode)
	Create_user(user *User_entity.User) global.ResponseStatusCode

	Update_user(e *User_entity.User) global.ResponseStatusCode
	Delete_user(id int, ctx context.Context) global.ResponseStatusCode
}

type User_service_impl struct {
	repo userrepo.User_repo
}

func (u *User_service_impl) Set_repo(r repo.Repo) {
	u.repo = r.(userrepo.User_repo)
	u.Init_key()
}
func (u *User_service_impl) Find_user_by_name(name string) (*User_entity.User, global.ResponseStatusCode) {
	e, err := u.repo.Get_by_name(name)
	if err != global.ResponseSuccess {
		return &User_entity.User{}, err
	} else {
		return e, global.ResponseSuccess
	}
}
func (u *User_service_impl) Find_user_by_id(id int) (*User_entity.User, global.ResponseStatusCode) {
	e, err := u.repo.Get_by_id(id)
	if err != global.ResponseSuccess {
		return &User_entity.User{}, err
	} else {
		return e, global.ResponseSuccess
	}
}
func (u *User_service_impl) hash_password(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes)
}
func (u *User_service_impl) Check_password(id int, password string) global.ResponseStatusCode {
	e, err := u.Find_user_by_id(id)
	if err != global.ResponseSuccess {
		return err
	}
	hash := e.Password
	ok := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if ok != nil {
		return global.ResponseIncorrectPassword
	}
	return global.ResponseSuccess
}

func (u *User_service_impl) Create_user(user *User_entity.User) global.ResponseStatusCode {
	user.Password = u.hash_password(user.Password)
	return u.repo.Create(user)
}

func (u *User_service_impl) Delete_user(id int, ctx context.Context) global.ResponseStatusCode {
	e := User_entity.User{
		Id:   id,
		Name: fmt.Sprintf("已注销用户_%d", id),
	}
	err := u.repo.ChangeUserNameByID(e.Id, e.Name)
	if err != global.ResponseSuccess {
		return err
	}
	err = u.repo.DestroyPassword(e.Id)
	if err != global.ResponseSuccess {
		return err
	}
	u.repo.UpdateActiveInRedisByUserId(e.Id, ctx)
	return global.ResponseSuccess

}
func (u *User_service_impl) Update_user(e *User_entity.User) global.ResponseStatusCode {
	if !(e.Password == "") {
		e.Password = u.hash_password(e.Password)
	}
	err := u.repo.Update(e)
	return err
}
