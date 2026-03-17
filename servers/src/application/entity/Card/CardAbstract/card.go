package CardAbstract

import (
	"pcc_card/application/entity/Card/character_card"
	"pcc_card/application/service/battleservice"
)

type Card interface {
	GetID() int
	GetInfo() map[string]any
}
type CardList struct {
	data []Card
}

var CardListImpl CardList

// InitCardList 对外接口
func InitCardList(s battleservice.BattleService) {
	CardListImpl = CardList{}
	CardListImpl.init(s)
}

func (c *CardList) Copy() *CardList {
	newList := &CardList{}
	if c.data == nil {
		return newList
	}
	newList.data = make([]Card, len(c.data))
	copy(newList.data, c.data)
	return newList
}

func (c *CardList) init(s battleservice.BattleService) {
	c.data = []Card{
		character_card.NewCard06(),
	}
}
