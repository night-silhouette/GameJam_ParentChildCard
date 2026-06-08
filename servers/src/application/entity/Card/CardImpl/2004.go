package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card2004 struct {
	CharacterBaseCard
}

func NewCard2004() *Card2004 {
	res := &Card2004{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2004) GetID() int {
	return 2004
}

func (c *Card2004) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card2004) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card2004) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
