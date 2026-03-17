package battleservice

import (
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/Card/CardImpl"
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
	for i := 0; i < len(c.data); i++ {
		= c.data[i]
	}
	return newList
}

func (c *CardList) init(s BattleService) {
	c.data = []CardAbstract.Card{
		CardImpl.NewCard00(),
		CardImpl.NewCard01(),
		CardImpl.NewCard02(),
		CardImpl.NewCard03(),
		CardImpl.NewCard04(),
		CardImpl.NewCard05(),
		CardImpl.NewCard06(),
		CardImpl.NewCard07(),
		CardImpl.NewCard08(),
		CardImpl.NewCard09(),
		CardImpl.NewCard10(),
		CardImpl.NewCard11(),
		CardImpl.NewCard12(),
		CardImpl.NewCard13(),
		CardImpl.NewCard14(),
		CardImpl.NewCard15(),
		CardImpl.NewCard16(),
		CardImpl.NewCard17(),
		CardImpl.NewCard18(),
		CardImpl.NewCard19(),
		CardImpl.NewCard20(),
		CardImpl.NewCard21(),
		CardImpl.NewCard22(),
		CardImpl.NewCard23(),
		CardImpl.NewCard24(),
		CardImpl.NewCard25(),
		CardImpl.NewCard26(),
		CardImpl.NewCard27(),
		CardImpl.NewCard28(),
		CardImpl.NewCard29(),
		CardImpl.NewCard30(),
		CardImpl.NewCard31(),
		CardImpl.NewCard32(),
		CardImpl.NewCard33(),
		CardImpl.NewCard34(),
		CardImpl.NewCard35(),
		CardImpl.NewCard36(),
		CardImpl.NewCard37(),
		CardImpl.NewCard38(),
		CardImpl.NewCard39(),
		CardImpl.NewCard40(),
		CardImpl.NewCard41(),
		CardImpl.NewCard42(),
		CardImpl.NewCard43(),
		CardImpl.NewCard44(),
	}
	for _, e := range c.data {
		info := s.GetCardInfoByID(e.GetID())
		e.SetInfo(info)
	}
}
