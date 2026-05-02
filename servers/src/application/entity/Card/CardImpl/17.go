package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card17 struct {
	BaseCard
}

func NewCard17() *Card17 {
	return &Card17{}
}

func (c *Card17) Attack(tempId int) {

}
func (c *Card17) Hurt(tempId int, HurtHp int) {
}

func (c *Card17) GetID() int {
	return 17
}

func (c *Card17) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card17) Skill(tempId int) {

}

func (c *Card17) Death(tempId int) {

}
