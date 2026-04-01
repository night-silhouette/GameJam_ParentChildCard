package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card04 struct {
	BaseCard
}

func NewCard04() *Card04 {
	return &Card04{}
}

func (c *Card04) GetID() int {
	return 4
}

func (c *Card04) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
