package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card0000 struct {
	CharacterBaseCard
}

func NewCard0000() *Card0000 {
	res := &Card0000{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0000) GetID() int {
	return 0
}

func (c *Card0000) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
