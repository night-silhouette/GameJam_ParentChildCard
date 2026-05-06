package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card05 struct {
	BaseCard
}

func NewCard05() *Card05 {
	return &Card05{}
}

func (c *Card05) GetID() int {
	return 5
}

func (c *Card05) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card05) PlayMagic() {}
