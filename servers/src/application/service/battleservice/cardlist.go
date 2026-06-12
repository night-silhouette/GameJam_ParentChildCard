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
	// 💡 不存实体了，改存“卡牌构造函数”的映射表，作为生产车间
	creators map[int]func() CardAbstract.Card
	// 💡 缓存配置表里的 is_parent 状态和基础 info，避免频繁跨服务查询
	cardInfoCache map[int]map[string]any
	Mt            sync.Mutex
	s             BattleService // 把服务实例挂载在工厂里，方便随时初始化新卡
}

var CardListImpl CardList

// InitCardList 对外全局初始化接口
func InitCardList(s BattleService) {
	CardListImpl = CardList{}
	CardListImpl.init(s)
}

// GetCardImpl 这是对外的、带锁的、安全的获取一张全新卡牌的方法
func (Cd *CardList) GetCardImpl(CardId int) CardAbstract.Card {
	Cd.Mt.Lock()
	defer Cd.Mt.Unlock()
	return Cd.getCardImpl(CardId)
}

// getCardImpl 内部核心工厂方法：动态 New 并当场完成所有动态变量注入
func (Cd *CardList) getCardImpl(CardId int) CardAbstract.Card {
	// 1. 通过 ID 匹配对应的构造器函数
	creator, exists := Cd.creators[CardId]
	if !exists {
		return nil
	}

	// 2. ⚡ 啪！直接从零 New 一张崭新、无污染的卡牌实体
	e := creator()

	// 3. 现场为这给新诞生的卡牌单独注入所有动态运行时数据
	e.InitBuffList()

	// 从缓存的图鉴配置中拿数据
	info := Cd.cardInfoCache[CardId]
	e.SetInfo(info)

	// ✨ 每个人都会拥有自己完全独立的、崭新的全局唯一 Channel，绝对不会多协程撞车打架！
	e.SetStateCodeChan(make(chan protocol.Effect))

	if val, ok := info["initHp"]; ok && val != nil {
		e.SetHpNow(info["initHp"].(float64))
	}
	if val, ok := info["damage"]; ok && val != nil {
		e.SetAtkNow(info["damage"].(float64))
	}

	// ✨ 分配全新独立的 Decorator 内存地址
	e.SetDec(CardMeta.NewDecorator())
	e.InitControlSignalMap()

	return e
}

// GetChildCard 输出所有子牌数组（全量生产崭新独立的子牌）
func (Cd *CardList) GetChildCard() []CardAbstract.Card {
	Cd.Mt.Lock()
	defer Cd.Mt.Unlock()

	res := make([]CardAbstract.Card, 0)
	// 遍历缓存，只要发现配置中 is_parent 为 false，说明是子牌，直接现场拉一条生产线
	for cardId, info := range Cd.cardInfoCache {
		if info != nil {
			if isParent, ok := info["is_parent"].(bool); ok && !isParent {
				// 直接调用 getCardImpl 生产出完全独立的实体
				if newChildCard := Cd.getCardImpl(cardId); newChildCard != nil {
					res = append(res, newChildCard)
				}
			}
		}
	}
	return res
}

// Copy 复制当前所有已注册卡牌的崭新列表（通常用于对局开始时初始化玩家的整套图鉴）
func (Cd *CardList) Copy() *[]CardAbstract.Card {
	Cd.Mt.Lock()
	defer Cd.Mt.Unlock()

	newList := make([]CardAbstract.Card, 0, len(Cd.creators))
	for cardId := range Cd.creators {
		if newCard := Cd.getCardImpl(cardId); newCard != nil {
			newList = append(newList, newCard)
		}
	}
	return &newList
}

// init 图鉴工厂初始化：在这里登記所有卡牌的出生证明
func (Cd *CardList) init(s BattleService) {
	Cd.s = s
	Cd.cardInfoCache = make(map[int]map[string]any)

	// 1. 注册所有的卡牌构造函数映射
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

	// 2. 提前把所有卡牌的静态 Info 配置拉取下来做成缓存，避免以后运行时反复请求 context 造成的延迟
	for cardId, creator := range Cd.creators {
		// 实例化一个临时卡牌用来拿它的 ID
		temp := creator()
		info := s.GetCardInfoByID(context.Background(), temp.GetID())
		Cd.cardInfoCache[cardId] = info
	}
}
