package BattleService

type Battle struct {
	BattleID int
	SM       StateMachine
	Ctx      Ctx
}

func (b *Battle) NewBattle(UserA int, UserB int) {

}
