package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card1007 struct {
	BaseCard
}

func NewCard1007() *Card1007 {
	return &Card1007{}
}

func (c *Card1007) PlayMagic() {}

func (c *Card1007) GetID() int {
	return 1007
}

func (c *Card1007) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
