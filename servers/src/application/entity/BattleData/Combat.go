package BattleData

type CombatDto struct {
	Behavior Behavior `json:"behavior" mapstructure:"behavior"`
	Target   Where    `json:"target" mapstructure:"target"`
}

func NewCombatDto(behavior Behavior, target Where) *CombatDto {
	res := CombatDto{}
	res.Behavior = behavior
	res.Target = target
	return &res
}

type Behavior int

const (
	Attack Behavior = iota
	Skill
)
