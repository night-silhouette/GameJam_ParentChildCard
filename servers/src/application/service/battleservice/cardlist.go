package battleservice

import (
	"context"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/Card/CardImpl"
	"sync"
)

type CardList struct {
	creators      map[int]func() CardAbstract.Card
	cardInfoCache map[int]map[string]any // 这里存的是“只读的原始配置图鉴”
	Mt            sync.Mutex
	s             BattleService
}

var CardListImpl CardList

func InitCardList(s BattleService) {
	CardListImpl = CardList{}
	CardListImpl.init(s)
}

// 根据tempid获取卡牌对象
func (Cd *CardList) GetCardImpl(CardId int) CardAbstract.Card {
	Cd.Mt.Lock()
	defer Cd.Mt.Unlock()
	return Cd.getCardImpl(CardId)
}

// getCardImpl 内部核心工厂方法（绝对安全版）
func (Cd *CardList) getCardImpl(CardId int) CardAbstract.Card {
	creator, exists := Cd.creators[CardId]
	if !exists {
		return nil
	}

	// 1. 生产独立新卡
	e := creator()

	// 2. 🛡️ 核心防污染安全区：把缓存的只读配置 Map，【深/浅拷贝】一份给新卡牌！
	cachedInfo := Cd.cardInfoCache[CardId]
	freshInfo := make(map[string]any) // 开辟完全属于这头新卡的独立 Map 内存
	for k, v := range cachedInfo {
		freshInfo[k] = v // 这样新卡牌随便怎么改自己的 freshInfo，全局缓存都绝对不会动！
	}

	// 把拷贝出来的独立 Info 塞给新卡
	e.SetInfo(freshInfo)

	// 4. 用 freshInfo 初始化数值
	if val, ok := freshInfo["initHp"]; ok && val != nil {
		e.SetHpNow(freshInfo["initHp"].(float64))
	}
	if val, ok := freshInfo["damage"]; ok && val != nil {
		e.SetAtkNow(freshInfo["damage"].(float64))
	}

	//------------------有新的初始化,就来这里------------------
	e.ShareInit()

	return e
}

// 获得子牌数组//这个里面也是用getimpl搞出来的独立对象
func (Cd *CardList) GetChildCard() []CardAbstract.Card {
	Cd.Mt.Lock()
	defer Cd.Mt.Unlock()

	res := make([]CardAbstract.Card, 0)
	for cardId, info := range Cd.cardInfoCache {
		if info != nil {
			if isParent, ok := info["is_parent"].(bool); ok && !isParent {
				if newChildCard := Cd.getCardImpl(cardId); newChildCard != nil {
					res = append(res, newChildCard)
				}
			}
		}
	}
	return res
}

func (Cd *CardList) init(s BattleService) {
	Cd.s = s
	Cd.cardInfoCache = make(map[int]map[string]any)

	Cd.creators = map[int]func() CardAbstract.Card{
		0: func() CardAbstract.Card { return CardImpl.NewCard0000() },
		1: func() CardAbstract.Card { return CardImpl.NewCard0001() },
		2: func() CardAbstract.Card { return CardImpl.NewCard0002() },
		3: func() CardAbstract.Card { return CardImpl.NewCard0003() },
		4: func() CardAbstract.Card { return CardImpl.NewCard0004() },
		5: func() CardAbstract.Card { return CardImpl.NewCard0005() },
		6: func() CardAbstract.Card { return CardImpl.NewCard0006() },
		7: func() CardAbstract.Card { return CardImpl.NewCard0007() },
		8: func() CardAbstract.Card { return CardImpl.NewCard0008() },
		9: func() CardAbstract.Card { return CardImpl.NewCard0009() },

		1000: func() CardAbstract.Card { return CardImpl.NewCard1000() },
		1001: func() CardAbstract.Card { return CardImpl.NewCard1001() },
		1002: func() CardAbstract.Card { return CardImpl.NewCard1002() },
		1003: func() CardAbstract.Card { return CardImpl.NewCard1003() },
		1004: func() CardAbstract.Card { return CardImpl.NewCard1004() },
		1005: func() CardAbstract.Card { return CardImpl.NewCard1005() },
		1006: func() CardAbstract.Card { return CardImpl.NewCard1006() },
		1007: func() CardAbstract.Card { return CardImpl.NewCard1007() },
		1008: func() CardAbstract.Card { return CardImpl.NewCard1008() },

		2000: func() CardAbstract.Card { return CardImpl.NewCard2000() },
		2001: func() CardAbstract.Card { return CardImpl.NewCard2001() },
		2002: func() CardAbstract.Card { return CardImpl.NewCard2002() },
		2003: func() CardAbstract.Card { return CardImpl.NewCard2003() },
		2004: func() CardAbstract.Card { return CardImpl.NewCard2004() },
		2005: func() CardAbstract.Card { return CardImpl.NewCard2005() },
		2006: func() CardAbstract.Card { return CardImpl.NewCard2006() },
		2007: func() CardAbstract.Card { return CardImpl.NewCard2007() },

		3000: func() CardAbstract.Card { return CardImpl.NewCard3000() },
		3001: func() CardAbstract.Card { return CardImpl.NewCard3001() },
		3002: func() CardAbstract.Card { return CardImpl.NewCard3002() },
		3003: func() CardAbstract.Card { return CardImpl.NewCard3003() },
		3004: func() CardAbstract.Card { return CardImpl.NewCard3004() },
		3005: func() CardAbstract.Card { return CardImpl.NewCard3005() },
		3006: func() CardAbstract.Card { return CardImpl.NewCard3006() },
		3007: func() CardAbstract.Card { return CardImpl.NewCard3007() },
	}

	// 缓存只读的原始图鉴数据
	for cardId, creator := range Cd.creators {
		temp := creator()
		info := s.GetCardInfoByID(context.Background(), temp.GetID())
		Cd.cardInfoCache[cardId] = info
	}
}
