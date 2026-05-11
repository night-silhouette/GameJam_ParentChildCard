package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card42 struct {
	CharacterBaseCard
}

func NewCard42() *Card42 {
	res := &Card42{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card42) GetID() int {
	return 42
}

func (c *Card42) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

