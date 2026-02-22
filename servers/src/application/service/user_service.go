package service

import "pcc_card/application/repo"

type User_service struct {
	repo repo.User_repo
}

var User_service_impl User_service

func NewUserService(repo repo.User_repo) {
	User_service_impl = User_service{repo}
}
