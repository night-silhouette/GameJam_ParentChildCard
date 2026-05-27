package battleservice

import (
	"context"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/service"
	"pcc_card/global"
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

	GiveCardByCardId(ctx context.Context, UserId int, CardId int) global.ResponseStatusCode        //给指定卡，没有保护，不可以用没有的cardId
	GiveInitCardBag(ctx context.Context, UserId int) global.ResponseStatusCode                     //给1级初始卡包
	GetBags(ctx context.Context, UserId int) ([]BattleData.BagStuffDto, global.ResponseStatusCode) //背包渲染
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
func (u *BattleServiceImpl) GiveCardByCardId(ctx context.Context, UserId int, CardId int) global.ResponseStatusCode {
	return u.repo.AddCardInBags(ctx, u.repo.Get_db(), CardId, UserId)
}

func (u *BattleServiceImpl) GiveInitCardBag(ctx context.Context, UserId int) global.ResponseStatusCode {
	CardIdList := make([]int, 0, 5)
	CardIdList = append(CardIdList, 0+Util.RandomRange(0, global.Lev1Category1Num))
	CardIdList = append(CardIdList, 0+Util.RandomRange(0, global.Lev1Category1Num))    //两张母战
	CardIdList = append(CardIdList, 1000+Util.RandomRange(0, global.Lev1Category2Num)) //一张母法
	CardIdList = append(CardIdList, 2000+Util.RandomRange(0, global.Lev1Category3Num)) //一张子战
	CardIdList = append(CardIdList, 3000+Util.RandomRange(0, global.Lev1Category4Num)) //一张子法
	tx, errDb := u.repo.Get_db().BeginTx(ctx, nil)
	if errDb != nil {
		return global.ResponseInternalServersError

	}
	defer tx.Rollback()
	for _, cardId := range CardIdList {
		u.repo.AddCardInBags(ctx, tx, cardId, UserId)
	}
	tx.Commit()
	return global.ResponseSuccess

}

func (u *BattleServiceImpl) GetBags(ctx context.Context, UserId int) ([]BattleData.BagStuffDto, global.ResponseStatusCode) {
	res, err := u.repo.GetBagsByUserId(ctx, u.repo.Get_db(), UserId)
	return res, err
}
