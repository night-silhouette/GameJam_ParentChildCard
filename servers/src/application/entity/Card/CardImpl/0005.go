package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card0005 struct {
	CharacterBaseCard
}

func NewCard0005() *Card0005 {
	res := &Card0005{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0005) GetID() int {
	return 5
}

func (c *Card0005) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
