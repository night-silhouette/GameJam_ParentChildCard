package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card3005 struct {
	BaseCard
}

func NewCard3005() *Card3005 {
	return &Card3005{}
}

func (c *Card3005) PlayMagic() {}

func (c *Card3005) GetID() int {
	return 3005
}

func (c *Card3005) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
