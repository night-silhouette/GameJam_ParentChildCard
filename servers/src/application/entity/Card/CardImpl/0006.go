package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card0006 struct {
	CharacterBaseCard
}

func NewCard0006() *Card0006 {
	res := &Card0006{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0006) GetID() int {
	return 6
}

func (c *Card0006) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
