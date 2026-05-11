package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card36 struct {
	CharacterBaseCard
}

func NewCard36() *Card36 {
	res := &Card36{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card36) GetID() int {
	return 36
}

func (c *Card36) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

