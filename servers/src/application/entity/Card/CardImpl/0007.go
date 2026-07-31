package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card0007 struct {
	CharacterBaseCard
	IsDie      bool //是不是碟
	SkillIsUse bool
}

func NewCard0007() *Card0007 {
	res := &Card0007{}
	res.CharacterBaseCard.Card = res
	res.IsDie = false
	res.SkillIsUse = false
	return res
}
func (c *Card0007) GetID() int {
	return 7
}

func (c *Card0007) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
func (c *Card0007) Skill(TargetId int) bool {
	FinalId := c.CheckGuard(TargetId)
	if !c.ShareSkill(FinalId) {
		return false
	}
	if c.SkillIsUse {
		return false
	}
	c.IsDie = true
	c.SkillIsUse = true
	c.Death(-1)
	return true
}

func (c *Card0007) Death(AttackId int) {
	if c.GetForm() == BattleData.NormalForm { //如果是普通形态,判定
		if c.IsDie {
			c.ChangeForm(BattleData.Die)
			return
		} else {
			c.ChangeForm(BattleData.E)
			return
		}
	}
	//优先级小于虫的变
	if c.BtCtx.GetWeather() == protocol.Fengdu && !c.changeJiangShi {
		c.ChangeForm(BattleData.JiangShi)
		c.changeJiangShi = true
		return //如果变僵尸了,就不用死了
	}

	c.Notify(BattleData.AnDeath, -1, AttackId, c.GetTempId())

	if !c.BtCtx.CheckIs2Bt(c.OwnerId) {
		config := protocol.NewInterruptConfig(c.OwnerId, c.BtCtx.ProtoColGetCharacterCard(c.OwnerId), 1, c.GetTempId(), BattleData.Deploy)
		c.BtCtx.ProtoColPush(protocol.NewSetCardBt(c.OwnerId, -1, true, &config))
		c.DisCard([]int{c.GetTempId()})
	}
}
