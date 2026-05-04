package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card10 struct {
	CharacterBaseCard
}

func NewCard10() *Card10 {
	return &Card10{}
}

func (c *Card10) GetID() int {
	return 10
}

func (c *Card10) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

