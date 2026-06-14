package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card0009 struct {
	CharacterBaseCard
}

func NewCard0009() *Card0009 {
	res := &Card0009{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0009) GetID() int {
	return 9
}

func (c *Card0009) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
