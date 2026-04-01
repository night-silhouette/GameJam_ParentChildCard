package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card29 struct {
	BaseCard
}

func NewCard29() *Card29 {
	return &Card29{}
}

func (c *Card29) Attack() {

}
func (c *Card29) Hurt() {
}

func (c *Card29) GetID() int {
	return 29
}

func (c *Card29) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
