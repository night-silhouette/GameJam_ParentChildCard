package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card29 struct {
	CharacterBaseCard
}

func NewCard29() *Card29 {
	res := &Card29{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card29) GetID() int {
	return 29
}

func (c *Card29) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

