package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card0008 struct {
	CharacterBaseCard
}

func NewCard0008() *Card0008 {
	res := &Card0008{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0008) GetID() int {
	return 8
}

func (c *Card0008) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
