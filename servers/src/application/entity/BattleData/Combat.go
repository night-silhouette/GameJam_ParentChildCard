package BattleData

type CombatDto struct {
	Behavior Behavior `json:"behavior" mapstructure:"behavior"`
	Where    Where    `json:"where" mapstructure:"where"`
}

func NewCombatDto(behavior Behavior, Where Where) *CombatDto {
	res := CombatDto{}
	res.Behavior = behavior
	res.Where = Where
	return &res
}

type Behavior int

const (
	Attack Behavior = iota
	Skill
)
