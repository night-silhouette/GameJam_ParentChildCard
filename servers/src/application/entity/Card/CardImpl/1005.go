package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card1005 struct {
	BaseCard
	TargetId int
}

func NewCard1005() *Card1005 {
	res := Card1005{}
	return &res
}

func (c *Card1005) PlayMagic() {
	config := protocol.NewInterruptConfig(c.OwnerId, c.BtCtx.ProtoColGetCharacterCard(c.OwnerId), 1, c.GetTempId(), BattleData.Selected)
	c.BtCtx.ProtoColPush(protocol.NewCustomInterrupt(func(res []int, pc protocol.ProtocolCardWithCtx) {
		c.TargetId = res[0]
	}, &config))
}

func (c *Card1005) GetID() int {
	return 1005
}

func (c *Card1005) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card1005) NextRound() {
	c.GiveBuff(&c.TargetId, *protocol.NewBuffBase(protocol.Untouchable, 1, 0, c.BtCtx.CreateTempId()))
}
