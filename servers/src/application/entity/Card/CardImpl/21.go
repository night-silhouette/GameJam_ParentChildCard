package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card21 struct {
	BaseCard
}

func NewCard21() *Card21 {
	return &Card21{}
}

func (c *Card21) Attack(tempId int) {

}
func (c *Card21) Hurt(tempId int, HurtHp int) {
}

func (c *Card21) GetID() int {
	return 21
}

func (c *Card21) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card21) Skill(tempId int) {

}

func (c *Card21) Death(tempId int) {

}
