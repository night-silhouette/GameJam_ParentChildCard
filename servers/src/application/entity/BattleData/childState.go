package BattleData

type ChildState int

const (
	Active ChildState = iota
	NotActive
	Died
	HasCatch
)
