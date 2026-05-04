package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card13 struct {
	CharacterBaseCard
}

func NewCard13() *Card13 {
	return &Card13{}
}

func (c *Card13) GetID() int {
	return 13
}

func (c *Card13) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

