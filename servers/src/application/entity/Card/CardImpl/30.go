package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card30 struct {
	CharacterBaseCard
}

func NewCard30() *Card30 {
	return &Card30{}
}

func (c *Card30) GetID() int {
	return 30
}

func (c *Card30) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

