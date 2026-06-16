package BattleData

//翻译卡牌数据，变成dto
//放这里的主要原因是因为包管理有点小问题，防止循环引用
//以后这里就放要返回前端的元数据

type CardDto struct {
	Id          int       `json:"id"`
	Hp          float64   `json:"hp"`
	Damage      float64   `json:"damage"`
	BuffDtoList []BuffDto `json:"buff_list"`
	TempId      int       `json:"temp_id"`
	Form        Form      `json:"form"`
}
