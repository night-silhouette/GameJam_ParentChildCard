package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card2006 struct {
	CharacterBaseCard
}

func NewCard2006() *Card2006 {
	res := &Card2006{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2006) GetID() int {
	return 2006
}

func (c *Card2006) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
