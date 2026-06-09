package BattleData

type DataAll struct {
	BtCardInfo    *BtCardInfo     `json:"bt_card_info" mapstructure:"bt_card_info"`
	CardInHand    *CardInHand     `json:"card_in_hand" mapstructure:"card_in_hand"`
	Energy        *EnergyDto      `json:"energy" mapstructure:"energy"`
	DiscardPool   []CardDto       `json:"discard_pool" mapstructure:"discard_pool"`
	ChildCardPool []*ChildCardDto `json:"child_child_pool" mapstructure:"child_child_pool"`
}

func NewDataAll(
	info *BtCardInfo,
	cardInHand *CardInHand,
	energy *EnergyDto,
	discardPool []CardDto,
	childCardPool []*ChildCardDto,
) *DataAll {
	return &DataAll{
		BtCardInfo:    info,
		CardInHand:    cardInHand,
		Energy:        energy,
		DiscardPool:   discardPool,
		ChildCardPool: childCardPool,
	}
}
