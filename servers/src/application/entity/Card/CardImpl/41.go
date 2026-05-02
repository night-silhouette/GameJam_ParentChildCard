package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card41 struct {
	BaseCard
}

func NewCard41() *Card41 {
	return &Card41{}
}

func (c *Card41) Attack(tempId int) {

}
func (c *Card41) Hurt(tempId int, HurtHp float64) {
}

func (c *Card41) GetID() int {
	return 41
}

func (c *Card41) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card41) Skill(tempId int) {

}

func (c *Card41) Death(tempId int) {

}
