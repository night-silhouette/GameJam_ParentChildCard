package db

import (
	"pcc_card/application/repo"
	"pcc_card/application/service"
)

func init_repo(repo repo.User_repo) {
	service.NewUserService(repo)
}
