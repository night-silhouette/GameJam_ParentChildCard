package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card35 struct {
	BaseCard
}

func NewCard35() *Card35 {
	return &Card35{}
}

func (c *Card35) GetID() int {
	return 35
}

func (c *Card35) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
