package battleservice

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/Card/character_card"
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
	InsertCard(data CardAbstract.Card)
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

func (u *BattleServiceImpl) InsertCard(id int) {
	c := character_card.CharacterCardTemplate{}
	c.ID = id
	c.Info["hp"] = 2
	c.Info["damage"] = 0
	u.repo.InsertCard(&c)
}
