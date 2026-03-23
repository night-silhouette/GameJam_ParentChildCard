package BattleData

type CardInHand struct {
	Self     []CardDto `json:"self"`
	Opponent []CardDto `json:"opponent"`
}

type CardDto struct {
	Id     int `json:"id"`
	Hp     int `json:"hp"`
	Damage int `json:"damage"`
	BuffId int `json:"buff_id"`
}
