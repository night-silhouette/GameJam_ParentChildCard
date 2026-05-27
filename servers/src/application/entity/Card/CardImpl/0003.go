package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card0003 struct {
	CharacterBaseCard
}

func NewCard0003() *Card0003 {
	res := &Card0003{}
	res.CharacterBaseCard.Card = res
	return res
}

func (c *Card0003) GetID() int {
	return 3
}

func (c *Card0003) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
