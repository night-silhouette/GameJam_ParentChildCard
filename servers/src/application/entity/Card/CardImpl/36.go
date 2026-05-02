package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card36 struct {
	BaseCard
}

func NewCard36() *Card36 {
	return &Card36{}
}

func (c *Card36) Attack(tempId int) {

}
func (c *Card36) Hurt(tempId int, HurtHp int) {
}

func (c *Card36) GetID() int {
	return 36
}

func (c *Card36) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card36) Skill(tempId int) {

}

func (c *Card36) Death(tempId int) {

}
