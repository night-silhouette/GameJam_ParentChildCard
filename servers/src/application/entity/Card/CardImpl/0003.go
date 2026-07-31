package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card0003 struct {
	CharacterBaseCard
}

func NewCard0003() *Card0003 {
	res := &Card0003{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0003) GetID() int {
	return 3
}
func (c *Card0003) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
func (c *Card0003) Skill(TargetId int) bool {
	FinalId := c.CheckGuard(TargetId)
	if !c.ShareSkill(FinalId) {
		return false
	}
	childCardId := c.BtCtx.GetTempIdByWhere(BattleData.ChildCard, c.GetOwnerId())
	c.GiveBuff(&childCardId, *protocol.NewBuffBase(protocol.Guard, 2, float64(childCardId), c.BtCtx.CreateTempId()))
	return true
}
