package battleservice

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/Card/character_card"
)

type CardList struct {
	data []CardAbstract.Card
}

var CardListImpl CardList

// InitCardList 对外接口
func InitCardList(s BattleService) {
	CardListImpl = CardList{}
	CardListImpl.init(s)
}

func (c *CardList) Copy() *CardList {
	newList := &CardList{}
	if c.data == nil {
		return newList
	}
	newList.data = make([]CardAbstract.Card, len(c.data))
	copy(newList.data, c.data)
	return newList
}

func (c *CardList) init(s BattleService) {
	c.data = []CardAbstract.Card{
		character_card.NewCard06(),
	}
}
