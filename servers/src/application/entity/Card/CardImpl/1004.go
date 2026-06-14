package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card1004 struct {
	BaseCard
}

func NewCard1004() *Card1004 {
	return &Card1004{}
}

func (c *Card1004) PlayMagic() {}

func (c *Card1004) GetID() int {
	return 1004
}

func (c *Card1004) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
