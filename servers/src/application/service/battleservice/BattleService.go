package battleservice

import (
	"context"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/service"
	"pcc_card/global"
	"pcc_card/infra/repo"
	"pcc_card/infra/repo/battlerepo"
	"pcc_card/infra/repo/userrepo"
	"sync"
)

type BattleService interface {
	service.Service
	AddMatch(id int, data BattleData.EnterBtData)
	IsHasID(id int) bool
	GetMatchSignals() *sync.Map
	GetCardInfoByID(ctx context.Context, ID int) map[string]any
	CheckUserIdIsBattle(ctx context.Context, userId int) (int, global.ResponseStatusCode)
}

type BattleServiceImpl struct {
	repo      battlerepo.BattleRepo
	User_repo userrepo.User_repo
}

func (u *BattleServiceImpl) Set_repo(r repo.Repo) { //注入对外接口
	u.repo = r.(battlerepo.BattleRepo)
}

func (u *BattleServiceImpl) AddMatch(id int, data BattleData.EnterBtData) {
	MM.AddPool(id, data)
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

func (u *BattleServiceImpl) CheckUserIdIsBattle(ctx context.Context, userId int) (int, global.ResponseStatusCode) {
	return u.User_repo.CheckUserIdIsBattle(ctx, u.repo.Get_db(), userId)
}
