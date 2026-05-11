package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card08 struct {
	CharacterBaseCard
}

func NewCard08() *Card08 {
	res := &Card08{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card08) GetID() int {
	return 8
}

func (c *Card08) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

