package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card23 struct {
	BaseCard
}

func NewCard23() *Card23 {
	return &Card23{}
}

func (c *Card23) GetID() int {
	return 23
}

func (c *Card23) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
