package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card21 struct {
	CharacterBaseCard
}

func NewCard21() *Card21 {
	res := &Card21{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card21) GetID() int {
	return 21
}

func (c *Card21) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

