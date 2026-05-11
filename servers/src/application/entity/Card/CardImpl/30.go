package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card30 struct {
	CharacterBaseCard
}

func NewCard30() *Card30 {
	res := &Card30{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card30) GetID() int {
	return 30
}

func (c *Card30) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

