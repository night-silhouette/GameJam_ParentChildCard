package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card02 struct {
	BaseCard
}

func NewCard02() *Card02 {
	return &Card02{}
}

func (c *Card02) GetID() int {
	return 2
}

func (c *Card02) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card02) PlayMagic() {}
