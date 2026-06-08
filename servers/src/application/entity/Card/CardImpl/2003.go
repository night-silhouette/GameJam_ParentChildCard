package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card2003 struct {
	CharacterBaseCard
}

func NewCard2003() *Card2003 {
	res := &Card2003{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2003) GetID() int {
	return 2003
}

func (c *Card2003) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card2003) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card2003) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
