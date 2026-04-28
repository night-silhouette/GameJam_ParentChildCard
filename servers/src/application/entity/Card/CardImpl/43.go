package CardImpl

import "pcc_card/application/entity/BattleData"


import "pcc_card/application/entity/Card/CardAbstract"


type Card43 struct {
	BaseCard
}

func NewCard43() *Card43 {
	return &Card43{}
}

func (c *Card43) Attack(w BattleData.Where) {

}
func (c *Card43) Hurt() {
}

func (c *Card43) GetID() int {
	return 43
}

func (c *Card43) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card43) Skill() {

}
