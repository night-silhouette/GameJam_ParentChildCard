package battleservice

import "pcc_card/presentation/handler/battlehandler/BattleDto"

type NotifyManager struct {
	ChanMap map[int]PlayerChannel // userid为键
}

type PlayerChannel struct {
	AcceptChan   chan BattleDto.Action
	ResponseChan chan BattleDto.Action
}

func NewNotifyManager(id1 int, id2 int, bufferSize int) *NotifyManager {
	nt := &NotifyManager{}
	nt.ChanMap = make(map[int]PlayerChannel, 2)
	nt.AddPlayer(id1, bufferSize)
	nt.AddPlayer(id2, bufferSize)
	return nt
}

func (nm *NotifyManager) AddPlayer(userID int, bufferSize int) {
	pc := PlayerChannel{
		AcceptChan:   make(chan BattleDto.Action, bufferSize),
		ResponseChan: make(chan BattleDto.Action, bufferSize),
	}
	nm.ChanMap[userID] = pc
}
