package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card19 struct {
	CharacterBaseCard
}

func NewCard19() *Card19 {
	res := &Card19{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card19) GetID() int {
	return 19
}

func (c *Card19) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

