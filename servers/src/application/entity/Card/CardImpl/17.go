package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card17 struct {
	CharacterBaseCard
}

func NewCard17() *Card17 {
	return &Card17{}
}

func (c *Card17) GetID() int {
	return 17
}

func (c *Card17) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

