package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card20 struct {
	CharacterBaseCard
}

func NewCard20() *Card20 {
	return &Card20{}
}

func (c *Card20) GetID() int {
	return 20
}

func (c *Card20) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

