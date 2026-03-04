package service

import (
	"database/sql"
	"errors"
	"fmt"
	"pcc_card/application/entity"
	"pcc_card/global"
	"pcc_card/infra/repo"
)

type User_service interface {
	Service
	Find_user_by_name(name string) (*entity.User, global.StatusCode)
}

type User_service_impl struct {
	repo repo.User_repo
}

func (u *User_service_impl) Set_repo(r repo.Repo) {
	u.repo = r.(repo.User_repo)
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
