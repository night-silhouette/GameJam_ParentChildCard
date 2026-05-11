package CardImpl

import (
	"pcc_card/application/entity/Card/CardAbstract"
)

type Card06 struct {
	CharacterBaseCard
}

func NewCard06() *Card06 {
	res := &Card06{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card06) GetID() int {
	return 6
}

func (c *Card06) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
