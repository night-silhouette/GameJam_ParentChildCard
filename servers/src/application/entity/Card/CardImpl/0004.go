package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card0004 struct {
	CharacterBaseCard
}

func NewCard0004() *Card0004 {
	res := &Card0004{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0004) GetID() int {
	return 4
}

func (c *Card0004) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
