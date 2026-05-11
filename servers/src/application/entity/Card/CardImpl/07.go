package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card07 struct {
	CharacterBaseCard
}

func NewCard07() *Card07 {
	res := &Card07{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card07) GetID() int {
	return 7
}

func (c *Card07) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
