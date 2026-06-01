package battleservice

import (
	"context"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/Card/CardImpl"
	"pcc_card/application/entity/CardMeta"
	"pcc_card/application/entity/protocol"
	"sync"
)

type CardList struct {
	data []CardAbstract.Card
	Mt   sync.Mutex
}

var CardListImpl CardList

// GetCardImpl 这是对外的,带锁的
func (Cd *CardList) GetCardImpl(CardId int) CardAbstract.Card {
	Cd.Mt.Lock()
	defer Cd.Mt.Unlock()
	return Cd.getCardImpl(CardId)
}

func (Cd *CardList) getCardImpl(CardId int) CardAbstract.Card {
	for _, el := range CardListImpl.data {
		if el.GetID() == CardId {
			res := el.Clone()
			return res
		}
	}
	return nil
}

// InitCardList 对外接口
func InitCardList(s BattleService) {
	CardListImpl = CardList{}
	CardListImpl.init(s)
}

// GetChildCard 输出子牌数组
func (Cd *CardList) GetChildCard() []CardAbstract.Card {
	Cd.Mt.Lock()
	defer Cd.Mt.Unlock()
	res := make([]CardAbstract.Card, 0)
	for _, el := range CardListImpl.data {
		if !el.GetInfo()["is_parent"].(bool) {
			res = append(res, el.Clone())
		}
	}
	return res
}

func (Cd *CardList) Copy() *[]CardAbstract.Card {
	var newList []CardAbstract.Card
	newList = make([]CardAbstract.Card, len(Cd.data))
	for i := 0; i < len(Cd.data); i++ {
		newList[i] = Cd.data[i].Clone()
	}
	return &newList
}

func (Cd *CardList) init(s BattleService) {
	Cd.data = []CardAbstract.Card{
		CardImpl.NewCard0000(),
		CardImpl.NewCard0001(),
		CardImpl.NewCard0002(),
		CardImpl.NewCard0003(),
		CardImpl.NewCard0004(),
		CardImpl.NewCard0005(),

		CardImpl.NewCard1000(),
		CardImpl.NewCard1001(),
		CardImpl.NewCard1002(),
		CardImpl.NewCard1003(),
		CardImpl.NewCard1004(),
		CardImpl.NewCard1005(),
		CardImpl.NewCard1006(),
		CardImpl.NewCard1007(),

		CardImpl.NewCard2000(),
		CardImpl.NewCard2001(),
		CardImpl.NewCard2002(),
		CardImpl.NewCard2003(),
		CardImpl.NewCard2004(),

		CardImpl.NewCard3000(),
		CardImpl.NewCard3001(),
		CardImpl.NewCard3002(),
		CardImpl.NewCard3003(),
		CardImpl.NewCard3004(),
	}
	for _, e := range Cd.data {
		e.InitBuffList()
		info := s.GetCardInfoByID(context.Background(), e.GetID())
		e.SetInfo(info)
		e.SetStateCodeChan(make(chan protocol.Effect))

		if val, ok := info["initHp"]; ok && val != nil {
			e.SetHpNow(info["initHp"].(float64))
		}
		if val, ok := info["damage"]; ok && val != nil {
			e.SetAtkNow(info["damage"].(float64))
		}

		e.SetDec(CardMeta.NewDecorator())
		e.InitControlSignalMap()
	}

}
