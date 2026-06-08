package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card2001 struct {
	CharacterBaseCard
}

func NewCard2001() *Card2001 {
	res := &Card2001{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2001) GetID() int {
	return 2001
}

func (c *Card2001) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card2001) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card2001) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
