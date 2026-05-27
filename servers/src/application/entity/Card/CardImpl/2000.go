package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card2000 struct {
	CharacterBaseCard
}

func NewCard2000() *Card2000 {
	res := &Card2000{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2000) GetID() int {
	return 2000
}

func (c *Card2000) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
