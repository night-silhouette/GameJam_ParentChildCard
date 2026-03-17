package battleservice

import (
	"sync"
	"sync/atomic"
)

var battleIDCounter int64

type Battle struct {
	mu       sync.RWMutex
	BattleID int
	SM       *StateMachine
	Ctx      *Ctx
}

func NewBattle(UserA int, UserB int) *Battle {
	id := int(atomic.AddInt64(&battleIDCounter, 1))
	ctx := Ctx{IDA: UserA, IDB: UserB}
	SM := NewStateMachine(&ctx)
	return &Battle{BattleID: id, SM: SM, Ctx: &ctx}
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

func (bc *BattleContainer) AddBattle(id1 int, id2 int) int {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	Bt := NewBattle(id1, id2)
	bc.Data[Bt.BattleID] = Bt
	bc.UserToBTID[id1] = Bt.BattleID
	bc.UserToBTID[id2] = Bt.BattleID
	return Bt.BattleID
}

func (bc *BattleContainer) GetBattle(id int) *Battle {
	bc.mu.RLock()
	defer bc.mu.RUnlock()
	BTID := bc.UserToBTID[id]
	BT := bc.Data[BTID]
	return BT
}
