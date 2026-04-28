package CardImpl

import "pcc_card/application/entity/BattleData"


import "pcc_card/application/entity/Card/CardAbstract"


type Card13 struct {
	BaseCard
}

func NewCard13() *Card13 {
	return &Card13{}
}

func (c *Card13) Attack(w BattleData.Where) {

}
func (c *Card13) Hurt() {
}

func (c *Card13) GetID() int {
	return 13
}

func (c *Card13) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card13) Skill() {

}
