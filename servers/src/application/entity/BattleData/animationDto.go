package BattleData

type AnimationBehavior int

const (
	AnAttack AnimationBehavior = iota
	AnHurt
	AnDeath
	AnSkill
)

type AnimationDto struct {
	Id                int               `json:"id" mapstructure:"id"`
	TempId            int               `json:"temp_id" mapstructure:"temp_id"`
	AnimationBehavior AnimationBehavior `json:"animation_behavior" mapstructure:"animation_behavior"`
	BtCardInfo        BtCardInfo        `json:"bt_card_info" mapstructure:"bt_card_info"`
}

func MewAnimationDto(id int, tempId int, AnimationBehavior AnimationBehavior, info BtCardInfo) AnimationDto {
	result := AnimationDto{}
	result.Id = id
	result.TempId = tempId
	result.AnimationBehavior = AnimationBehavior
	result.BtCardInfo = info
	return result
}
