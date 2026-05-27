package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card2004 struct {
	CharacterBaseCard
}

func NewCard2004() *Card2004 {
	res := &Card2004{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2004) GetID() int {
	return 2004
}

func (c *Card2004) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
