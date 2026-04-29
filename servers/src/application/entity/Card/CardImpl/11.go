package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card11 struct {
	BaseCard
}

func NewCard11() *Card11 {
	return &Card11{}
}

func (c *Card11) Attack(tempId int) {

}
func (c *Card11) Hurt(tempId int) {
}

func (c *Card11) GetID() int {
	return 11
}

func (c *Card11) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card11) Skill(tempId int) {

}
