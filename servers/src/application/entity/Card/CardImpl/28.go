package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card28 struct {
	CharacterBaseCard
}

func NewCard28() *Card28 {
	res := &Card28{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card28) GetID() int {
	return 28
}

func (c *Card28) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

