package battleservice

import (
	"context"
	"fmt"
	"pcc_card/application/entity/Card/CardAbstract"
	"pcc_card/application/entity/protocol"
	"pcc_card/infra/repo/userrepo"
	"sync"
	"sync/atomic"
)

var battleIDCounter int64

type Battle struct {
	mu sync.RWMutex //房间锁

	Context context.Context //生命周期总控制
	Cancel  context.CancelFunc

	BattleID int
	SM       *StateMachine
	Ctx      *Ctx
	Nt       *NotifyManager

	TempId atomic.Int32
}

// 传来的卡的id。出来的是充血对象(填充tempid了的),tempid是自增了的
func CloneCardListByid(cardIdList []int, NumCalc *atomic.Int32, GoCtx context.Context, ctx protocol.ProtocolCardWithCtx, CtxRecord *CardAbstract.CtxRecord) []CardAbstract.Card {
	ChildCardList := make([]CardAbstract.Card, 0)

	for _, CardId := range cardIdList {
		e := CardListImpl.GetCardImpl(CardId, GoCtx, ctx, CtxRecord)
		e.SetTempId(int(NumCalc.Add(1)))
		ChildCardList = append(ChildCardList, e)
	}
	return ChildCardList
}

func NewBattle(UserA int, UserB int, CardList map[int][]int, GoldMoreUserId int, cList []int) *Battle {
	var TempId atomic.Int32
	TempId.Store(int32(0))
	rootContext := context.Background()
	BattleContext, cancel := context.WithCancel(rootContext)
	id := int(atomic.AddInt64(&battleIDCounter, 1))
	CtxRecord := CardAbstract.NewCtxRecord()
	c := Ctx{}
	//clone手牌,给userid到ownerId
	//这个函数递增了tempid的计数,

	CardInHand := make(map[int]map[int]CardAbstract.Card)
	for UserId, Value := range CardList {
		CardInHand[UserId] = make(map[int]CardAbstract.Card)
		List := CloneCardListByid(Value, &TempId, BattleContext, &c, CtxRecord)
		for _, Card := range List {
			CardInHand[UserId][Card.GetTempId()] = Card
		}
	}
	//这里就设置了手牌的主人,子牌的主人为没有,获取的时候设置
	for UserId, CardMap := range CardInHand { //卡的主人这种信息属于局内信息,让ctx和battle自己管,而注入ctx,goctx和各种底层初始化,要解偶开
		for _, Card := range CardMap {
			Card.SetOwnerId(UserId)
		}
	}
	//全都聚合在这里的,匹配只是传过来id,他只指挥生成什么牌,不具体生成牌,因为ctx和GoCtx都在这里
	//子牌堆初始化
	ChildCardList := make([]CardAbstract.Card, 0)
	ChildCardList = CloneCardListByid(cList, &TempId, BattleContext, &c, CtxRecord)

	InitCtx(&c, UserA, UserB, BattleContext, CardInHand, &TempId, ChildCardList, CtxRecord)
	Nt := NewNotifyManager(UserA, UserB, 32) //初始化bufferSize
	SM := NewStateMachine(&c, UserA, UserB, Nt, BattleContext, GoldMoreUserId)
	go func() {
		select {
		case <-BattleContext.Done():
			BC.RemoveBattle(id)
		}
	}()
	return &Battle{BattleID: id, SM: SM, Ctx: &c, Nt: Nt, Context: BattleContext, Cancel: cancel}
}

func (b *Battle) GetPlayerChanByUserID(id int) PlayerChannel {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.Nt.ChanMap[id]
}

//------------------------------------------------------------

type BattleContainer struct {
	mu         sync.RWMutex
	Data       map[int]*Battle
	UserToBTID map[int]int
	User_repo  userrepo.User_repo
}

var BC BattleContainer

func (b *BattleContainer) GetBattleData() map[int]*Battle {
	b.mu.RLock()
	defer b.mu.RUnlock()
	fmt.Println("GetBattleContainer", b.Data)
	return b.Data
}

func InitBattleContainer(repo userrepo.User_repo) {
	BC = BattleContainer{}
	BC.User_repo = repo
	BC.Data = make(map[int]*Battle)
	BC.UserToBTID = make(map[int]int)
	battleIDCounter = 1

}

// AddBattle 传来的卡的id。
func (bc *BattleContainer) AddBattle(id1 int, id2 int, cardIdList map[int][]int, GoldMoreUserId int, cList []int) int { //启动接口

	Bt := NewBattle(id1, id2, cardIdList, GoldMoreUserId, cList)
	bc.mu.Lock()
	defer bc.mu.Unlock()
	bc.Data[Bt.BattleID] = Bt
	bc.UserToBTID[id1] = Bt.BattleID
	bc.UserToBTID[id2] = Bt.BattleID

	return Bt.BattleID
}

func (bc *BattleContainer) GetBattleByUserID(id int) *Battle {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	BTID := bc.UserToBTID[id]
	if BTID == 0 {
		return nil
	}
	BT := bc.Data[BTID]
	return BT
}
func (bc *BattleContainer) RemoveBattle(BattleId int) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	delete(bc.Data, BattleId)
	delete(bc.UserToBTID, BattleId)
}
