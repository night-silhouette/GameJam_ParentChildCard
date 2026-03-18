package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card24 struct {
	BaseCard
}

func NewCard24() *Card24 {
	return &Card24{}
}

func (c *Card24) GetID() int {
	return 24
}

func (c *Card24) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
