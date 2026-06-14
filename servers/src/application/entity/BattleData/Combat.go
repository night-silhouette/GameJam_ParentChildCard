package BattleData

type CombatDto struct {
	Behavior      Behavior   `json:"behavior" mapstructure:"behavior"`
	SelfWhere     Where      `json:"self_where" mapstructure:"self_where"`
	OpponentWhere Where      `json:"opponent_where" mapstructure:"opponent_where"`
	SelectCard    SelectCard `json:"select_card" mapstructure:"select_card"`
}

func NewCombatDto(behavior Behavior, SelfWhere Where, OpponentWhere Where) *CombatDto {
	res := CombatDto{}
	res.Behavior = behavior
	res.SelfWhere = SelfWhere
	res.OpponentWhere = OpponentWhere
	return &res
}

type Behavior int

const (
	Attack Behavior = iota
	Skill
	SwitchCard
)
