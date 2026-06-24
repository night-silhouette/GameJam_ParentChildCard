package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card1002 struct {
	BaseCard
}

func NewCard1002() *Card1002 {
	return &Card1002{}
}

func (c *Card1002) PlayMagic() {
	config := protocol.NewInterruptConfig(c.OwnerId, c.BtCtx.ProtoColGetCharacterCard(c.OwnerId), 2, c.GetTempId(), BattleData.Selected)
	c.BtCtx.ProtoColPush(protocol.NewCustomInterrupt(func(res []int, pc protocol.ProtocolCardWithCtx) {
		A := res[0]
		B := res[1]
		pc.ProtoCol1002(A, B)
	}, &config))
}

func (c *Card1002) GetID() int {
	return 1002
}

func (c *Card1002) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
