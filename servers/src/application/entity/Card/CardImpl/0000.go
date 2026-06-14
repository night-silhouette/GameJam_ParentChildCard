package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card0000 struct {
	CharacterBaseCard
	SKillMapList []map[int]int //map是 buffTempId:num
}

func NewCard0000() *Card0000 {
	res := &Card0000{}
	res.CharacterBaseCard.Card = res
	res.SKillMapList = make([]map[int]int, 0)
	return res
}

func (c *Card0000) GetID() int {
	return 0
}

func (c *Card0000) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card0000) Skill(TargetId int) bool {
	if !c.CharacterBaseCard.Skill(TargetId) {
		return false
	}
	TempId := c.GetTempId()
	BuffTempId := c.BtCtx.CreateTempId()
	c.GiveBuff(&TempId, *protocol.NewBuffBase(protocol.DamageImmunity, 3, 0.5, BuffTempId)) //这个和直接上buff没什么区别
	c.SKillMapList = append(c.SKillMapList)
	Map := make(map[int]int)
	Map[BuffTempId] = 2 //使用次数
	c.SKillMapList = append(c.SKillMapList, Map)
	return true
}

func (c *Card0000) Hurt(AttackId int, HurtHp float64, category BattleData.ValueChange) {
	c.NewCustom(func(pc protocol.ProtocolCardWithCtx) {
		for _, Map := range c.SKillMapList { //受到伤害,减少使用次数
			for BuffTempId, Num := range Map {
				Num -= 1
				if Num == 0 {
					c.ReMoveBuffByTempId(BuffTempId) //这个方法安全的,随便删,即使时间到了,也删
				}
			}
		}
	})
	c.CharacterBaseCard.Hurt(AttackId, HurtHp, category) //反向塞进去,先扣血,再减少层数
}
