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
}

func MewAnimationDto(id int, tempId int, AnimationBehavior AnimationBehavior) AnimationDto {
	result := AnimationDto{}
	result.Id = id
	result.TempId = tempId
	result.AnimationBehavior = AnimationBehavior
	return result
}
