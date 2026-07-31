package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card1001 struct {
	BaseCard
}

func NewCard1001() *Card1001 {
	return &Card1001{}
}

func (c *Card1001) PlayMagic() {
	config := protocol.NewInterruptConfig(c.OwnerId, c.BtCtx.ProtoColGetCharacterCard(c.BtCtx.GetOpponentId(c.OwnerId)), 1, c.GetTempId(), BattleData.Selected)
	c.BtCtx.ProtoColPush(protocol.NewGiveBuff(&c.TempId, *protocol.NewBuffBase(protocol.HealingDecay, 3, -2, c.BtCtx.CreateTempId()), true, &config))
}

func (c *Card1001) GetID() int {
	return 1001
}

func (c *Card1001) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
