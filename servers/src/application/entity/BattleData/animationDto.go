package BattleData

type AnimationBehavior int

const (
	AnAttack AnimationBehavior = iota
	AnHurt
	AnDeath
	AnSkill
	AnDisCard
)

type AnimationDto struct {
	TempId            int               `json:"temp_id" mapstructure:"temp_id"`
	AnimationBehavior AnimationBehavior `json:"animation_behavior" mapstructure:"animation_behavior"`
	BtCardInfo        BtCardInfo        `json:"bt_card_info" mapstructure:"bt_card_info"`
}

func NewAnimationDto(tempId int, AnimationBehavior AnimationBehavior, info BtCardInfo) AnimationDto {
	result := AnimationDto{}
	result.TempId = tempId
	result.AnimationBehavior = AnimationBehavior
	result.BtCardInfo = info
	return result
}
