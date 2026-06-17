package BattleData

type FormChangeDto struct {
	Form    Form    `json:"form"` //传的是改变后的form
	TempId  int     `json:"temp_id"`
	DataAll DataAll `json:"data_all"`
}

func NewFormChangeDto(form Form, TempId int, all DataAll) *FormChangeDto {
	result := new(FormChangeDto)
	result.Form = form
	result.TempId = TempId
	result.DataAll = all
	return result
}
