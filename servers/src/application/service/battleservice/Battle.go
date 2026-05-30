package battleservice

import (
	"context"
	"pcc_card/application/entity/Card/CardAbstract"
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

func NewBattle(UserA int, UserB int, CardList map[int][]int, GoldMoreUserId int) *Battle {
	var TempId atomic.Int32
	TempId.Store(int32(0))
	rootContext := context.Background()
	BattleContext, cancel := context.WithCancel(rootContext)
	id := int(atomic.AddInt64(&battleIDCounter, 1))
	ctx := NewCtx(UserA, UserB, CardListImpl.Copy(), BattleContext, CloneByCardListImpl(CardList, &TempId), &TempId)
	Nt := NewNotifyManager(UserA, UserB, 32) //初始化bufferSize
	SM := NewStateMachine(ctx, UserA, UserB, Nt, BattleContext, GoldMoreUserId)
	go func() {
		select {
		case <-BattleContext.Done():
			BC.RemoveBattle(id)
		}
	}()
	return &Battle{BattleID: id, SM: SM, Ctx: ctx, Nt: Nt, Context: BattleContext, Cancel: cancel}
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
	return b.Data
}

func InitBattleContainer(repo userrepo.User_repo) {
	BC = BattleContainer{}
	BC.User_repo = repo
	BC.Data = make(map[int]*Battle)
	BC.UserToBTID = make(map[int]int)
	battleIDCounter = 1

}

// CloneByCardListImpl 传来的卡的id。出来的是充血对象(填充tempid了的)
func CloneByCardListImpl(cardIdList map[int][]int, NumCalc *atomic.Int32) map[int]map[int]CardAbstract.Card {
	res := make(map[int]map[int]CardAbstract.Card)
	for key, value := range cardIdList {
		res[key] = make(map[int]CardAbstract.Card)
		for _, CardId := range value {
			card := CardListImpl.GetCardImpl(CardId)

			tempId := (*NumCalc).Add(1)
			card.SetTempId(int(tempId))
			res[key][card.GetTempId()] = card
		}

	}
	return res
}

// AddBattle 传来的卡的id。
func (bc *BattleContainer) AddBattle(id1 int, id2 int, cardIdList map[int][]int, GoldMoreUserId int) int { //启动接口

	Bt := NewBattle(id1, id2, cardIdList, GoldMoreUserId)
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
