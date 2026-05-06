package BattleData

type BtCardInfo struct {
	Self     BtCardInfoMeta `json:"self" mapstructure:"self"`
	Opponent BtCardInfoMeta `json:"opponent" mapstructure:"opponent"`
}

type BtCardInfoMeta struct {
	SkillCardBt  CardDto `json:"skill_card_bt" mapstructure:"skill_card_bt"`
	ParentCardBt CardDto `json:"parent_card_bt" mapstructure:"parent_card_bt"`
	ChildCardBt  CardDto `json:"child_card_bt" mapstructure:"child_card_bt"`
}
