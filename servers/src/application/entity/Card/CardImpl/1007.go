package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card1007 struct {
	BaseCard
}

func NewCard1007() *Card1007 {
	return &Card1007{}
}

func (c *Card1007) PlayMagic() {
	c.EffectUpdateEnergy(2)
	c.NewCustom(func(pc protocol.ProtocolCardWithCtx) {
		pc.ChangeWeather(protocol.Ganmu)
	})
}

func (c *Card1007) GetID() int {
	return 1007
}

func (c *Card1007) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
