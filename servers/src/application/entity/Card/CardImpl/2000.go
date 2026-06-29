package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card2000 struct {
	CharacterBaseCard
	HurtId int //上一次的hurttempid//上一次被激发的时候,是被那次攻击激发,同一次激发,不会触发
}

func NewCard2000() *Card2000 {
	res := &Card2000{}
	res.CharacterBaseCard.Card = res
	res.HurtId = 0
	return res
}

func (c *Card2000) GetID() int {
	return 2000
}

func (c *Card2000) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card2000) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) {
		if c.HurtId != c.CtxRecord.LastHurt.TempId && c.CtxRecord.LastHurt.Value >= 6 {
			UserId := pc.CardUserIdByTempId(c.CtxRecord.LastHurt.TempId)
			c.HurtId = c.CtxRecord.LastHurt.TempId
			return true, UserId
		}
		return false, 0
	})
}

func (c *Card2000) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
