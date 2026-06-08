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
	Caller            int               `json:"caller" mapstructure:"caller"`
	Acceptor          int               `json:"acceptor" mapstructure:"acceptor"`
	AnimationBehavior AnimationBehavior `json:"animation_behavior" mapstructure:"animation_behavior"`
	DataAll           DataAll           `json:"data_all" mapstructure:"data_all"`
}

type DataAll struct {
	BtCardInfo    BtCardInfo     `json:"bt_card_info" mapstructure:"bt_card_info"`
	CardInHand    CardInHand     `json:"card_in_hand" mapstructure:"card_in_hand"`
	Energy        int            `json:"energy" mapstructure:"energy"`
	DiscardPool   []CardDto      `json:"discard_pool" mapstructure:"discard_pool"`
	ChildCardPool []ChildCardDto `json:"child_child_pool" mapstructure:"child_child_pool"`
}

func NewAnimationDto(
	callerId int,
	acceptorId int,
	behavior AnimationBehavior,
	info BtCardInfo,
	cardInHand CardInHand,
	energy int,
	discardPool []CardDto,
	childCardPool []ChildCardDto,
) AnimationDto {
	return AnimationDto{
		Caller:            callerId,
		Acceptor:          acceptorId,
		AnimationBehavior: behavior,
		DataAll: DataAll{
			BtCardInfo:    info,
			CardInHand:    cardInHand,
			Energy:        energy,
			DiscardPool:   discardPool,
			ChildCardPool: childCardPool,
		},
	}
}
