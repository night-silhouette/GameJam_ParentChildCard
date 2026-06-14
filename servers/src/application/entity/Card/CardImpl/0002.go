package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card0002 struct {
	CharacterBaseCard
}

func NewCard0002() *Card0002 {
	res := &Card0002{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0002) GetID() int {
	return 2
}

func (c *Card0002) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
