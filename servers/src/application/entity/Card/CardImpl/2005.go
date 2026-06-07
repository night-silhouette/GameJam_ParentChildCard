package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card2005 struct {
	CharacterBaseCard
}

func NewCard2005() *Card2005 {
	res := &Card2005{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2005) GetID() int {
	return 2005
}

func (c *Card2005) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
