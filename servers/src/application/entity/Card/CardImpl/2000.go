package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card2000 struct {
	CharacterBaseCard
}

func NewCard2000() *Card2000 {
	res := &Card2000{}
	res.CharacterBaseCard.Card = res
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
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card2000) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
