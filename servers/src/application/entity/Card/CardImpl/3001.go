package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)



type Card3001 struct {
	BaseCard
}

func NewCard3001() *Card3001 {
	return &Card3001{}
}

func (c *Card3001) PlayMagic() {}

func (c *Card3001) GetID() int {
	return 3001
}

func (c *Card3001) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card3001) Check(pc protocol.ProtocolCardWithCtx) CardAbstract.ChildCheckFunc {
	return CardAbstract.ChildCheckFunc(func(pc protocol.ProtocolCardWithCtx) (bool, int) { return false, 0 })
}

func (c *Card3001) Trigger(pc protocol.ProtocolCardWithCtx, UserId int) {

}
