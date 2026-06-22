package BattleData

type SelectCard struct {
	Where      Where `json:"where" mapstructure:"where"`
	CardId     int   `json:"card_id" mapstructure:"card_id"`
	CardTempId int   `json:"card_temp_id" mapstructure:"card_temp_id"`
}

//下牌之后,是否还有战斗牌
//

type Where int

const (
	ParentCard Where = iota
	ChildCard
	SkillCard
	DisCardPool
	ChildCardPool
	InHand
	Self
)
