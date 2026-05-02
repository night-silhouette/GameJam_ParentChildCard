package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card10 struct {
	BaseCard
}

func NewCard10() *Card10 {
	return &Card10{}
}

func (c *Card10) Attack(tempId int) {

}
func (c *Card10) Hurt(tempId int, HurtHp float64) {
}

func (c *Card10) GetID() int {
	return 10
}

func (c *Card10) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card10) Skill(tempId int) {

}

func (c *Card10) Death(tempId int) {

}
