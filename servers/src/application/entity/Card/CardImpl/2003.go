package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card2003 struct {
	CharacterBaseCard
}

func NewCard2003() *Card2003 {
	res := &Card2003{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2003) GetID() int {
	return 2003
}

func (c *Card2003) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
