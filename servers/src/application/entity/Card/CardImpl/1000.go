package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card1008 struct {
	BaseCard
}

func NewCard1008() *Card1008 {
	return &Card1008{}
}

func (c *Card1008) PlayMagic() {}

func (c *Card1008) GetID() int {
	return 1008
}

func (c *Card1008) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
