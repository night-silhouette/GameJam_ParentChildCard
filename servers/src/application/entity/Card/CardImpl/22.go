package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card22 struct {
	BaseCard
	CharacterBaseCard
}

func NewCard22() *Card22 {
	return &Card22{}
}

func (c *Card22) GetID() int {
	return 22
}

func (c *Card22) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

