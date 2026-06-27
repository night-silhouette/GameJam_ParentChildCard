package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card1003 struct {
	BaseCard
}

func NewCard1003() *Card1003 {
	return &Card1003{}
}

func (c *Card1003) PlayMagic() {
	config := protocol.NewInterruptConfig(c.OwnerId, c.BtCtx.ProtoColGetCharacterCard(c.OwnerId), 1, c.GetTempId(), BattleData.Selected)
	c.BtCtx.ProtoColPush(protocol.NewGiveBuff(&c.TempId, *protocol.NewBuffBase(protocol.DamageTransform, 3, 0.4, c.BtCtx.CreateTempId()), true, &config))
}

func (c *Card1003) GetID() int {
	return 1003
}

func (c *Card1003) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
