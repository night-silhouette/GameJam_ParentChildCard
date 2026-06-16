package BattleData

type CardCalcValueDto struct {
	TempId   int         `json:"temp_id" mapstructure:"temp_id"`
	Category ValueChange `json:"category" mapstructure:"category"`
	Value    float64     `json:"value" mapstructure:"value"` //有正有负
	DataAll  *DataAll    `json:"data_all" mapstructure:"data_all"`
	IsMiss   bool        `json:"is_miss" mapstructure:"is_miss"`
}

type ValueChange int

const (
	Damage ValueChange = iota
	Heal
	TrueDamage
)
