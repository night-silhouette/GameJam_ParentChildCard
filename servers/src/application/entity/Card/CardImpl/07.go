package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card07 struct {
	BaseCard
}

func NewCard07() *Card07 {
	return &Card07{}
}

func (c *Card07) Attack() {

}
func (c *Card07) Hurt() {
}

func (c *Card07) GetID() int {
	return 7
}

func (c *Card07) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
