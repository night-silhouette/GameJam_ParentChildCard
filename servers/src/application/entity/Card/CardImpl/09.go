package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card09 struct {
	BaseCard
}

func NewCard09() *Card09 {
	return &Card09{}
}

func (c *Card09) Attack() {

}
func (c *Card09) Hurt() {
}

func (c *Card09) GetID() int {
	return 9
}

func (c *Card09) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card09) Skill() {

}
