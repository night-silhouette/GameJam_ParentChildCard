package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card14 struct {
	BaseCard
}

func NewCard14() *Card14 {
	return &Card14{}
}

func (c *Card14) Attack() {

}
func (c *Card14) Hurt() {
}

func (c *Card14) GetID() int {
	return 14
}

func (c *Card14) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
