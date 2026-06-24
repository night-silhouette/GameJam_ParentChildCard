package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
)

type Card1006 struct {
	BaseCard
}

func NewCard1006() *Card1006 {
	return &Card1006{}
}

func (c *Card1006) PlayMagic() {
	config := protocol.NewInterruptConfig(c.OwnerId, c.BtCtx.ProtoColGetCharacterCard(c.BtCtx.GetOpponentId(c.OwnerId)), 1, c.GetTempId(), BattleData.Selected)
	c.BtCtx.ProtoColPush(protocol.NewGiveBuff(&c.TempId, *protocol.NewBuffBase(protocol.Confine, 4, 0, c.BtCtx.CreateTempId()), true, &config))
}

func (c *Card1006) GetID() int {
	return 1006
}

func (c *Card1006) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
