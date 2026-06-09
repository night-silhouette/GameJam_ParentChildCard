package BattleData

type EnergyDto struct {
	Self     int `json:"self" mapstructure:"self"`
	Opponent int `json:"opponent" mapstructure:"opponent"`
}
