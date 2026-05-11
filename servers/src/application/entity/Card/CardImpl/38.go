package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card38 struct {
	CharacterBaseCard
}

func NewCard38() *Card38 {
	res := &Card38{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card38) GetID() int {
	return 38
}

func (c *Card38) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

