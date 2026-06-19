package BattleData

import (
	"pcc_card/Util"
	"time"
)

type InterruptDto struct {
	StateWaitTime int64         `json:"state_wait_time" mapstructure:"state_wait_time"`
	TempIdList    []int         `json:"temp_id_list" mapstructure:"temp_id_list"`
	SelectNum     int           `json:"select_num" mapstructure:"select_num"`
	InterruptType InterruptType `json:"interrupt_type" mapstructure:"interrupt_type"`
	CallTempId    int           `json:"call_temp_id" mapstructure:"call_temp_id"`
}

type InterruptSelect struct {
	TempIdList []int `json:"temp_id_list" mapstructure:"temp_id_list"`
}

func NewInterruptDto(StateWaitTime time.Duration, TempIdList []int, SelectNum int, InterruptType InterruptType, CallTempId int) *InterruptDto {
	res := InterruptDto{}
	res.StateWaitTime = Util.SendTime(StateWaitTime)
	res.TempIdList = TempIdList
	res.SelectNum = SelectNum
	return &res
}

type InterruptType int

const (
	Selected InterruptType = iota
	Deploy
)
