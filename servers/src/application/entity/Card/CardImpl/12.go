package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card12 struct {
	CharacterBaseCard
}

func NewCard12() *Card12 {
	res := &Card12{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card12) GetID() int {
	return 12
}

func (c *Card12) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

