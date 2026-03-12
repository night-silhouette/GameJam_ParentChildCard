package battleservice

import (
	"sync"
	"sync/atomic"
)

var battleIDCounter int64

type Battle struct {
	BattleID int
	SM       StateMachine
	Ctx      Ctx
}

func NewBattle(UserA int, UserB int) *Battle {
	id := int(atomic.AddInt64(&battleIDCounter, 1))
	Ctx := Ctx{UserA, UserB, make(map[string]any)}
	SM := NewStateMachine()
	return &Battle{BattleID: id, SM: *SM, Ctx: Ctx}
}

//------------------------------------------------------------

type BattleContainer struct {
	mu   sync.RWMutex
	Data map[int]*Battle
}

var BC BattleContainer

func InitBattleContainer() {
	BC = BattleContainer{}
	BC.Data = make(map[int]*Battle)
}

func (bc *BattleContainer) AddBattle(id1 int, id2 int) {
	bc.mu.Lock()
	defer bc.mu.Unlock()
	NewBattle(id1, id2)
}
