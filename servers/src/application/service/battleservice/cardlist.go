package battleservice

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/Card/character_card"
	"pcc_card/application/entity/Card/skill_card"
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
		skill_card.NewCard00(),
		skill_card.NewCard01(),
		skill_card.NewCard02(),
		skill_card.NewCard03(),
		skill_card.NewCard04(),
		skill_card.NewCard05(),
		character_card.NewCard06(),
		character_card.NewCard07(),
		character_card.NewCard08(),
		character_card.NewCard09(),
		character_card.NewCard10(),
		character_card.NewCard11(),
		character_card.NewCard12(),
		character_card.NewCard13(),
		character_card.NewCard14(),
		character_card.NewCard15(),
		character_card.NewCard16(),
		character_card.NewCard17(),
		character_card.NewCard18(),
		character_card.NewCard19(),
		character_card.NewCard20(),
		character_card.NewCard21(),
		character_card.NewCard22(),
		skill_card.NewCard23(),
		skill_card.NewCard24(),
		skill_card.NewCard25(),
		character_card.NewCard26(),
		character_card.NewCard27(),
		character_card.NewCard28(),
		character_card.NewCard29(),
		character_card.NewCard30(),
		character_card.NewCard31(),
		skill_card.NewCard32(),
		skill_card.NewCard33(),
		skill_card.NewCard34(),
		skill_card.NewCard35(),
		character_card.NewCard36(),
		character_card.NewCard37(),
		character_card.NewCard38(),
		character_card.NewCard39(),
		character_card.NewCard40(),
		character_card.NewCard41(),
		character_card.NewCard42(),
		character_card.NewCard43(),
		skill_card.NewCard44(),
	}
	for _, e := range c.data {
		info := s.GetCardInfoByID(e.GetID())
		e.SetInfo(info)
	}
}
