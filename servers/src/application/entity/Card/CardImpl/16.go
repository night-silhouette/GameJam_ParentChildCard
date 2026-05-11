package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card16 struct {
	CharacterBaseCard
}

func NewCard16() *Card16 {
	res := &Card16{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card16) GetID() int {
	return 16
}

func (c *Card16) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

