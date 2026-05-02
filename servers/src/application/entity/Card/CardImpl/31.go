package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card31 struct {
	BaseCard
	CharacterBaseCard
}

func NewCard31() *Card31 {
	return &Card31{}
}

func (c *Card31) GetID() int {
	return 31
}

func (c *Card31) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

