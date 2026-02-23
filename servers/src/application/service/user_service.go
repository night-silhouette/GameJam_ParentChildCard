package service

import (
	"pcc_card/application/repo"
)

type User_service struct {
	repo repo.User_repo
}

func NewUserService(repo repo.User_repo) *User_service {
	return &User_service{repo}
}
