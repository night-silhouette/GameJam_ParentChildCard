package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card3002 struct {
	BaseCard
}

func NewCard3002() *Card3002 {
	return &Card3002{}
}

func (c *Card3002) PlayMagic() {}

func (c *Card3002) GetID() int {
	return 3002
}

func (c *Card3002) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
