package battleservice

import (
	"sync"
	"sync/atomic"
)

var battleIDCounter int64

type Battle struct {
	mu       sync.RWMutex //房间锁
	BattleID int
	SM       *StateMachine
	Ctx      *Ctx
	Nt       *NotifyManager
}

func NewBattle(UserA int, UserB int) *Battle {
	id := int(atomic.AddInt64(&battleIDCounter, 1))
	ctx := NewCtx(UserA, UserB)
	Nt := NewNotifyManager(UserA, UserB, 32) //初始化bufferSize
	SM := NewStateMachine(ctx, UserA, UserB, Nt)
	return &Battle{BattleID: id, SM: SM, Ctx: ctx, Nt: Nt}
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

func InitBattleContainer() {
	BC = BattleContainer{}
	BC.Data = make(map[int]*Battle)
	BC.UserToBTID = make(map[int]int)
}

func (bc *BattleContainer) AddBattle(id1 int, id2 int) int { //启动接口
	bc.mu.Lock()
	defer bc.mu.Unlock()
	Bt := NewBattle(id1, id2)
	bc.Data[Bt.BattleID] = Bt
	bc.UserToBTID[id1] = Bt.BattleID
	bc.UserToBTID[id2] = Bt.BattleID
	return Bt.BattleID
}

func (bc *BattleContainer) GetBattleByUserID(id int) *Battle {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	BTID := bc.UserToBTID[id]
	BT := bc.Data[BTID]
	return BT
}
