package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card2007 struct {
	CharacterBaseCard
}

func NewCard2007() *Card2007 {
	res := &Card2007{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2007) GetID() int {
	return 2007
}

func (c *Card2007) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
