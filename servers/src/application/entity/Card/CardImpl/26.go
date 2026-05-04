package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card26 struct {
	CharacterBaseCard
}

func NewCard26() *Card26 {
	return &Card26{}
}

func (c *Card26) GetID() int {
	return 26
}

func (c *Card26) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

