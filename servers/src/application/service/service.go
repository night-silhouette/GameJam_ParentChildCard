package service

import "pcc_card/infra/repo"

type Service interface {
	Set_repo(repo *repo.Repo)
}
