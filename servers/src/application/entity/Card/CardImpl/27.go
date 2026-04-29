package CardImpl

import "pcc_card/application/entity/Card/CardAbstract"


type Card27 struct {
	BaseCard
}

func NewCard27() *Card27 {
	return &Card27{}
}

func (c *Card27) Attack(tempId int) {

}
func (c *Card27) Hurt(tempId int) {
}

func (c *Card27) GetID() int {
	return 27
}

func (c *Card27) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card27) Skill(tempId int) {

}
