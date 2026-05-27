package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card1005 struct {
	BaseCard
}

func NewCard1005() *Card1005 {
	return &Card1005{}
}

func (c *Card1005) PlayMagic() {}

func (c *Card1005) GetID() int {
	return 1005
}

func (c *Card1005) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
