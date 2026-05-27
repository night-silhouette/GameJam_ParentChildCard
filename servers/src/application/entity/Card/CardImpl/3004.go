package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card3004 struct {
	BaseCard
}

func NewCard3004() *Card3004 {
	return &Card3004{}
}

func (c *Card3004) PlayMagic() {}

func (c *Card3004) GetID() int {
	return 3004
}

func (c *Card3004) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
