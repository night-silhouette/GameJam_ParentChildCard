package BattleData

type CardInHand struct {
	Self     []CardDto `json:"self"`
	Opponent []CardDto `json:"opponent"`
}
