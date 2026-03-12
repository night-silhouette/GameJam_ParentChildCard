package BattleService

import (
	"pcc_card/application/service"
	"pcc_card/infra/repo"
	"pcc_card/infra/repo/battle_repo"
)

type BattleService interface {
	service.Service
}

type battleServiceimpl struct {
	repo battle_repo.BattleRepo
}

func (u *battleServiceimpl) Set_repo(r repo.Repo) { //注入对外接口
	u.repo = r.(battle_repo.BattleRepo)
}
