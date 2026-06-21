package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card3006 struct {
	BaseCard
}

func NewCard3006() *Card3006 {
	return &Card3006{}
}

func (c *Card3006) PlayMagic() {}

func (c *Card3006) GetID() int {
	return 3006
}

func (c *Card3006) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card3006) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card3006) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
