package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card18 struct {
	CharacterBaseCard
}

func NewCard18() *Card18 {
	res := &Card18{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card18) GetID() int {
	return 18
}

func (c *Card18) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

