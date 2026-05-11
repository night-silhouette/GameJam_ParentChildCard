package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card43 struct {
	CharacterBaseCard
}

func NewCard43() *Card43 {
	res := &Card43{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card43) GetID() int {
	return 43
}

func (c *Card43) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

