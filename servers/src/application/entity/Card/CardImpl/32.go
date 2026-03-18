package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card32 struct {
	BaseCard
}

func NewCard32() *Card32 {
	return &Card32{}
}

func (c *Card32) GetID() int {
	return 32
}

func (c *Card32) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
