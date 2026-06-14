package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card1001 struct {
	BaseCard
}

func NewCard1001() *Card1001 {
	return &Card1001{}
}

func (c *Card1001) PlayMagic() {}

func (c *Card1001) GetID() int {
	return 1001
}

func (c *Card1001) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
