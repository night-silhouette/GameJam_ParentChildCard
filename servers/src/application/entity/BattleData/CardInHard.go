package BattleData

type CardInHand struct {
	Self     []CardDto `json:"self"`
	Opponent []CardDto `json:"opponent"`
}

type CardDto struct {
	Id     int     `json:"id"`
	Hp     float64 `json:"hp"`
	Damage float64 `json:"damage"`
	BuffId int     `json:"buff_id"`
}
