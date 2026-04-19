package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card01 struct {
	BaseCard
}

func (c *Card01) PlayMagic() {}

func NewCard01() *Card01 {
	return &Card01{}
}

func (c *Card01) GetID() int {
	return 1
}

func (c *Card01) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
