package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card3001 struct {
	BaseCard
}

func NewCard3001() *Card3001 {
	return &Card3001{}
}

func (c *Card3001) PlayMagic() {}

func (c *Card3001) GetID() int {
	return 3001
}

func (c *Card3001) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
