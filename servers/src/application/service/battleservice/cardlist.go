package battleservice

import (
	"context"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/Card/CardImpl"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/application/entity/protocol"
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

func (c *CardList) Copy() *[]CardAbstract.Card {
	var newList []CardAbstract.Card
	newList = make([]CardAbstract.Card, len(c.data))
	for i := 0; i < len(c.data); i++ {
		newList[i] = c.data[i].Clone()
	}
	return &newList
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
		e.InitBuffList()
		info := s.GetCardInfoByID(context.Background(), e.GetID())
		e.SetInfo(info)
		e.SetStateCodeChan(make(chan protocol.Effect))

		if val, ok := info["hp"]; ok && val != nil {
			e.SetHpNow(info["hp"].(float64))
		}
		if val, ok := info["damage"]; ok && val != nil {
			e.SetAtkNow(info["damage"].(float64))
		}
		if e.GetID() == 15 {
			e.SetHpNow(3)
		}
		e.SetDec(CardMeta.NewDecorator())
		e.InitControlSignalMap()
	}

}
