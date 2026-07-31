package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card1008 struct {
	BaseCard
}

func NewCard1008() *Card1008 {
	return &Card1008{}
}

func (c *Card1008) PlayMagic() {
	config := protocol.NewInterruptConfig(c.OwnerId, c.BtCtx.ProtoColGetCharacterCard(c.OwnerId), 1, c.GetTempId(), BattleData.Selected)
	c.BtCtx.ProtoColPush(protocol.NewGiveBuff(&c.TempId, *protocol.NewBuffBase(protocol.BonusDamage, 4, 2, c.BtCtx.CreateTempId()), true, &config))
}

func (c *Card1008) GetID() int {
	return 1008
}

func (c *Card1008) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
