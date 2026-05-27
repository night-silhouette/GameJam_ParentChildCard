package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card2001 struct {
	CharacterBaseCard
}

func NewCard2001() *Card2001 {
	res := &Card2001{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2001) GetID() int {
	return 2001
}

func (c *Card2001) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
