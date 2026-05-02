package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card31 struct {
	BaseCard
}

func NewCard31() *Card31 {
	return &Card31{}
}

func (c *Card31) Attack(tempId int) {

}
func (c *Card31) Hurt(tempId int, HurtHp int) {
}

func (c *Card31) GetID() int {
	return 31
}

func (c *Card31) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card31) Skill(tempId int) {

}

func (c *Card31) Death(tempId int) {

}
