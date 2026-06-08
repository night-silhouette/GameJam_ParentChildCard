package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card3007 struct {
	BaseCard
}

func NewCard3007() *Card3007 {
	return &Card3007{}
}

func (c *Card3007) PlayMagic() {}

func (c *Card3007) GetID() int {
	return 3007
}

func (c *Card3007) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card3007) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card3007) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
