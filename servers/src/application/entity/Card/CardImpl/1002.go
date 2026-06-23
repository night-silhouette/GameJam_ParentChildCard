package CardImpl

import (
	"pcc_card/application/entity/BattleData"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/global"
	"time"
)

type Card1002 struct {
	BaseCard
}

func NewCard1002() *Card1002 {
	return &Card1002{}
}

func (c *Card1002) PlayMagic() {
	CheckIsInterrupt := true
	InterruptSelect := make([]int, 0)

	ToSelect := c.BtCtx.ProtoColGetCharacterCard(c.BtCtx.GetOpponentId(c.OwnerId))
	c.Interrupt(&InterruptSelect, global.Interrupt*time.Second, ToSelect, 1, &CheckIsInterrupt, BattleData.Selected)
	c.GiveBuff()
}

func (c *Card1002) GetID() int {
	return 1002
}

func (c *Card1002) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
