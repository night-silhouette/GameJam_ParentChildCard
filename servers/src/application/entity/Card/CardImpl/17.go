package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card17 struct {
	CharacterBaseCard
}

func NewCard17() *Card17 {
	res := &Card17{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card17) GetID() int {
	return 17
}

func (c *Card17) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

