package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card2002 struct {
	CharacterBaseCard
}

func NewCard2002() *Card2002 {
	res := &Card2002{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card2002) GetID() int {
	return 2002
}

func (c *Card2002) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
