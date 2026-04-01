package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card38 struct {
	BaseCard
}

func NewCard38() *Card38 {
	return &Card38{}
}

func (c *Card38) Attack() {

}
func (c *Card38) Hurt() {
}

func (c *Card38) GetID() int {
	return 38
}

func (c *Card38) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}
