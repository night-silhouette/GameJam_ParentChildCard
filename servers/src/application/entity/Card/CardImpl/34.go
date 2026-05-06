package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card34 struct {
	BaseCard
}

func NewCard34() *Card34 {
	return &Card34{}
}

func (c *Card34) GetID() int {
	return 34
}

func (c *Card34) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
func (c *Card34) PlayMagic() {}
