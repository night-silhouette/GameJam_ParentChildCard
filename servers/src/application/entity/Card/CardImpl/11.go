package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card11 struct {
	CharacterBaseCard
}

func NewCard11() *Card11 {
	res := &Card11{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card11) GetID() int {
	return 11
}

func (c *Card11) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

