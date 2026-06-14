package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card3003 struct {
	BaseCard
}

func NewCard3003() *Card3003 {
	return &Card3003{}
}

func (c *Card3003) PlayMagic() {}

func (c *Card3003) GetID() int {
	return 3003
}

func (c *Card3003) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card3003) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card3003) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
