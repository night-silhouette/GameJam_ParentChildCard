package service

import (
	"pcc_card/infra/repo"
)

type user_service interface {
	Service
}

type User_service_impl struct {
	repo repo.Repo
}

func (u *User_service_impl) Set_repo(repo repo.Repo) {
	u.repo = repo
}
