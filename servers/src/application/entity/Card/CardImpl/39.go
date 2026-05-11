package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card39 struct {
	CharacterBaseCard
}

func NewCard39() *Card39 {
	res := &Card39{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card39) GetID() int {
	return 39
}

func (c *Card39) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

