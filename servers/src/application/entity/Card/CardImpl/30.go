package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card30 struct {
	BaseCard
}

func NewCard30() *Card30 {
	return &Card30{}
}

func (c *Card30) Attack(tempId int) {

}
func (c *Card30) Hurt(tempId int, HurtHp float64) {
}

func (c *Card30) GetID() int {
	return 30
}

func (c *Card30) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card30) Skill(tempId int) {

}

func (c *Card30) Death(tempId int) {

}
