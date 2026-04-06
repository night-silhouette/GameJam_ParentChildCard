package battleservice

import (
	"context"
	"pcc_card/application/service"
	"pcc_card/infra/repo"
	"pcc_card/infra/repo/battlerepo"
	"sync"
)

type BattleService interface {
	service.Service
	AddMatch(id int)
	IsHasID(id int) bool
	GetMatchSignals() *sync.Map
	GetCardInfoByID(ctx context.Context, ID int) map[string]any
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
func (u *BattleServiceImpl) GetMatchSignals() *sync.Map {
	return &MatchSignals
}

func (u *BattleServiceImpl) GetCardInfoByID(ctx context.Context, ID int) map[string]any {
	return u.repo.ReadCardByID(ctx, u.repo.Get_db(), ID)
}
