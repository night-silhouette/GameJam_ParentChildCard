package CardImpl

import "pcc_card/application/entity/BattleData"


import "pcc_card/application/entity/Card/CardAbstract"


type Card15 struct {
	BaseCard
}

func NewCard15() *Card15 {
	return &Card15{}
}

func (c *Card15) Attack(w BattleData.Where) {

}
func (c *Card15) Hurt() {
}

func (c *Card15) GetID() int {
	return 15
}

func (c *Card15) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card15) Skill() {

}
