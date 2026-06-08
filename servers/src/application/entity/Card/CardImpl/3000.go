package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card3000 struct {
	BaseCard
}

func NewCard3000() *Card3000 {
	return &Card3000{}
}

func (c *Card3000) PlayMagic() {}

func (c *Card3000) GetID() int {
	return 3000
}

func (c *Card3000) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card3000) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card3000) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
