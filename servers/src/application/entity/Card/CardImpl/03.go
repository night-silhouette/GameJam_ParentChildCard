package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"

type Card03 struct {
	BaseCard
}

func NewCard03() *Card03 {
	return &Card03{}
}

func (c *Card03) GetID() int {
	return 3
}

func (c *Card03) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
func (c *Card03) PlayMagic() {}
