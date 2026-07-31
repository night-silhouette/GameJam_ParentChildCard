package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card1004 struct {
	BaseCard
}

func NewCard1004() *Card1004 {
	return &Card1004{}
}

func (c *Card1004) PlayMagic() {
	List := c.BtCtx.ProtoGetBtAll(c.GetOwnerId())
	for _, v := range List {
		c.EffectHeal(v, 5) //全体出战者治疗5血
	}
	c.NewCustom(func(pc protocol.ProtocolCardWithCtx) {
		pc.ChangeWeather(protocol.Ganmu)
	})
}

func (c *Card1004) GetID() int {
	return 1004
}

func (c *Card1004) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
