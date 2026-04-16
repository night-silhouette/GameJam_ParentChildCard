package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card25 struct {
	BaseCard
}

func NewCard25() *Card25 {
	return &Card25{}
}

func (c *Card25) GetID() int {
	return 25
}

func (c *Card25) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
func (c *Card25) PlayMagic() {}
