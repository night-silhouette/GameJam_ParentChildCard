package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card39 struct {
	BaseCard
}

func NewCard39() *Card39 {
	return &Card39{}
}

func (c *Card39) Attack(tempId int) {

}
func (c *Card39) Hurt(tempId int, HurtHp float64) {
}

func (c *Card39) GetID() int {
	return 39
}

func (c *Card39) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card39) Skill(tempId int) {

}

func (c *Card39) Death(tempId int) {

}
