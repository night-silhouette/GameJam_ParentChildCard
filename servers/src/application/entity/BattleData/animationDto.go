package BattleData

type AnimationBehavior int

const (
	AnAttack AnimationBehavior = iota
	AnHurt
	AnDeath
	AnSkill
)

type AnimationDto struct {
	Caller            int               `json:"caller" mapstructure:"caller"`
	Acceptor          int               `json:"acceptor" mapstructure:"acceptor"`
	AnimationBehavior AnimationBehavior `json:"animation_behavior" mapstructure:"animation_behavior"`
}

func NewAnimationDto(callerId int, acceptorId int, behavior AnimationBehavior) AnimationDto {
	return AnimationDto{
		Caller:            callerId,
		Acceptor:          acceptorId,
		AnimationBehavior: behavior,
	}
}
