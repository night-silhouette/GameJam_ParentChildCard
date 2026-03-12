package BattleService

import "sync/atomic"

var battleIDCounter int64

type Battle struct {
	BattleID int
	SM       StateMachine
	Ctx      Ctx
}

func (b *Battle) NewBattle(UserA int, UserB int) *Battle {
	id := int(atomic.AddInt64(&battleIDCounter, 1))
	Ctx := Ctx{UserA, UserB, make(map[string]any)}
	SM := NewStateMachine()
	return &Battle{BattleID: id, SM: *SM, Ctx: Ctx}
}
