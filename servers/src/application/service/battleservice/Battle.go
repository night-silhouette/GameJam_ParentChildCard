package battleservice

import (
	"context"
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
}

func NewBattle(UserA int, UserB int) *Battle {
	rootContext := context.Background()
	BattleContext, cancel := context.WithCancel(rootContext)
	id := int(atomic.AddInt64(&battleIDCounter, 1))
	ctx := NewCtx(UserA, UserB, CardListImpl.Copy(), BattleContext)
	Nt := NewNotifyManager(UserA, UserB, 32) //初始化bufferSize
	SM := NewStateMachine(ctx, UserA, UserB, Nt, BattleContext)
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
}

var BC BattleContainer

func (b *BattleContainer) GetBattleData() map[int]*Battle {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.Data
}

func InitBattleContainer() {
	BC = BattleContainer{}
	BC.Data = make(map[int]*Battle)
	BC.UserToBTID = make(map[int]int)
	battleIDCounter = 1

}

func (bc *BattleContainer) AddBattle(id1 int, id2 int) int { //启动接口
	Bt := NewBattle(id1, id2)
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
