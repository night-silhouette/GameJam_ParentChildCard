package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card1006 struct {
	BaseCard
}

func NewCard1006() *Card1006 {
	return &Card1006{}
}

func (c *Card1006) PlayMagic() {}

func (c *Card1006) GetID() int {
	return 1006
}

func (c *Card1006) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
