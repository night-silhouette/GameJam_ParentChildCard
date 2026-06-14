package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card2002 struct {
	CharacterBaseCard
}

func NewCard2002() *Card2002 {
	res := &Card2002{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2002) GetID() int {
	return 2002
}

func (c *Card2002) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card2002) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card2002) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
