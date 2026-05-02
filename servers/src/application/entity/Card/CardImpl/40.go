package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card40 struct {
	BaseCard
	CharacterBaseCard
}

func NewCard40() *Card40 {
	return &Card40{}
}

func (c *Card40) GetID() int {
	return 40
}

func (c *Card40) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

