package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card41 struct {
	CharacterBaseCard
}

func NewCard41() *Card41 {
	res := &Card41{}
    res.CharacterBaseCard.Card=res
    return res
}

func (c *Card41) GetID() int {
	return 41
}

func (c *Card41) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

