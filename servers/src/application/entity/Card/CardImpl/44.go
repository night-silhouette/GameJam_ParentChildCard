package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card44 struct {
	BaseCard
}

func NewCard44() *Card44 {
	return &Card44{}
}

func (c *Card44) GetID() int {
	return 44
}

func (c *Card44) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
