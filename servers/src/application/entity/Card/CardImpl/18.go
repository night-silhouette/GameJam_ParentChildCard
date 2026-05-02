package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card18 struct {
	BaseCard
}

func NewCard18() *Card18 {
	return &Card18{}
}

func (c *Card18) Attack(tempId int) {

}
func (c *Card18) Hurt(tempId int, HurtHp int) {
}

func (c *Card18) GetID() int {
	return 18
}

func (c *Card18) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card18) Skill(tempId int) {

}

func (c *Card18) Death(tempId int) {

}
