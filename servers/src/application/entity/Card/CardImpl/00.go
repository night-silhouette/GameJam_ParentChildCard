package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card00 struct {
	BaseCard
}

func NewCard00() *Card00 {
	return &Card00{}
}

func (c *Card00) PlayMagic() {}

func (c *Card00) GetID() int {
	return 0
}

func (c *Card00) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
