package battleservice

type NotifyManager struct {
	ChanMap map[int]PlayerChannel // userid为键
}

type PlayerChannel struct {
	AcceptChan   chan any
	ResponseChan chan any
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
		AcceptChan:   make(chan any, bufferSize),
		ResponseChan: make(chan any, bufferSize),
	}
	nm.ChanMap[userID] = pc
}
