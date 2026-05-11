package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card22 struct {
	CharacterBaseCard
}

func NewCard22() *Card22 {
	res := &Card22{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card22) GetID() int {
	return 22
}

func (c *Card22) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

