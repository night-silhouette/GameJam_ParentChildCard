package BattleData

type CardCalcCardMoveDto struct {
	Object  Where    `json:"object" mapstructure:"object"`
	TempId  int      `json:"temp_id" mapstructure:"temp_id"`
	DataAll *DataAll `json:"data_all" mapstructure:"data_all"`
}
