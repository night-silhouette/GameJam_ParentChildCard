package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card2006 struct {
	CharacterBaseCard
}

func NewCard2006() *Card2006 {
	res := &Card2006{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2006) GetID() int {
	return 2006
}

func (c *Card2006) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card2006) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card2006) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
