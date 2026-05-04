package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card08 struct {
	CharacterBaseCard
}

func NewCard08() *Card08 {
	return &Card08{}
}

func (c *Card08) GetID() int {
	return 8
}

func (c *Card08) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

