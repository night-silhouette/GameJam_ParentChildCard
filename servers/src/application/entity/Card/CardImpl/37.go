package CardImpl

import "pcc_card/application/entity/BattleData"


import "pcc_card/application/entity/Card/CardAbstract"


type Card37 struct {
	BaseCard
}

func NewCard37() *Card37 {
	return &Card37{}
}

func (c *Card37) Attack(w BattleData.Where) {

}
func (c *Card37) Hurt() {
}

func (c *Card37) GetID() int {
	return 37
}

func (c *Card37) Clone() CardAbstract.Card {
	newCard := *c
	return &newCard
}

func (c *Card37) Skill() {

}
