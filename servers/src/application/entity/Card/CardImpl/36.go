package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card36 struct {
	BaseCard
	CharacterBaseCard
}

func NewCard36() *Card36 {
	return &Card36{}
}

func (c *Card36) GetID() int {
	return 36
}

func (c *Card36) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

