package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card1000 struct {
	BaseCard
}

func NewCard1000() *Card1000 {
	return &Card1000{}
}

func (c *Card1000) PlayMagic() {}

func (c *Card1000) GetID() int {
	return 1000
}

func (c *Card1000) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
