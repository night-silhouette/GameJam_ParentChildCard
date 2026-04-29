package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card22 struct {
	BaseCard
}

func NewCard22() *Card22 {
	return &Card22{}
}

func (c *Card22) Attack(tempId int) {

}
func (c *Card22) Hurt(tempId int) {
}

func (c *Card22) GetID() int {
	return 22
}

func (c *Card22) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card22) Skill(tempId int) {

}
