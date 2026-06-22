package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card1000 struct {
	BaseCard
}

func NewCard1000() *Card1000 {
	return &Card1000{}
}

func (c *Card1000) PlayMagic() {
	cList := c.BtCtx.GetCharacterId()
	for _, id := range cList {
		c.BtCtx.ProtoColPush(protocol.NewCustom(func(pc protocol.ProtocolCardWithCtx) {
			pc.ProtoColAttackNoHurt(id, -1, BattleData.Damage)
		}))
	}
}

func (c *Card1000) GetID() int {
	return 1000
}

func (c *Card1000) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
