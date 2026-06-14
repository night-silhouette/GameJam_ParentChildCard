package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card0007 struct {
	CharacterBaseCard
}

func NewCard0007() *Card0007 {
	res := &Card0007{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0007) GetID() int {
	return 7
}

func (c *Card0007) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
