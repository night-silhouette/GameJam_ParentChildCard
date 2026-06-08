package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card3002 struct {
	BaseCard
}

func NewCard3002() *Card3002 {
	return &Card3002{}
}

func (c *Card3002) PlayMagic() {}

func (c *Card3002) GetID() int {
	return 3002
}

func (c *Card3002) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card3002) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card3002) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
