package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card15 struct {
	CharacterBaseCard
}

func NewCard15() *Card15 {
	res := &Card15{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card15) GetID() int {
	return 15
}

func (c *Card15) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

