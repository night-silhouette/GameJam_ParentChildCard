package service

import "pcc_card/application/repo"

type User_service struct {
	repo repo.User_repo
}

var UserService User_service

func NewUserService(repo repo.User_repo) {
	UserService = User_service{repo}
}
