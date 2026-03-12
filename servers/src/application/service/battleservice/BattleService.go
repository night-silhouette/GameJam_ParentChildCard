package battleservice

import (
	"pcc_card/application/service"
	"pcc_card/infra/repo"
	"pcc_card/infra/repo/battlerepo"
)

type BattleService interface {
	service.Service
	AddMatch(id int)
	IsHasID(id int) bool
}

type BattleServiceImpl struct {
	repo battlerepo.BattleRepo
}

func (u *BattleServiceImpl) Set_repo(r repo.Repo) { //注入对外接口
	u.repo = r.(battlerepo.BattleRepo)
}

func (u *BattleServiceImpl) AddMatch(id int) {
	MM.AddPool(id)
}
func (u *BattleServiceImpl) IsHasID(id int) bool {
	return MM.IsHasID(id)
}
