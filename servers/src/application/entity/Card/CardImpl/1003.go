package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card1003 struct {
	BaseCard
}

func NewCard1003() *Card1003 {
	return &Card1003{}
}

func (c *Card1003) PlayMagic() {}

func (c *Card1003) GetID() int {
	return 1003
}

func (c *Card1003) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
