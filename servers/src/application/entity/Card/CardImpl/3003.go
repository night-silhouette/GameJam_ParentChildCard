package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card3003 struct {
	BaseCard
}

func NewCard3003() *Card3003 {
	return &Card3003{}
}

func (c *Card3003) PlayMagic() {}

func (c *Card3003) GetID() int {
	return 3003
}

func (c *Card3003) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
