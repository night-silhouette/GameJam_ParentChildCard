package BattleService

import (
	"pcc_card/global"
	"sync"
	"time"
)

type MatchManager struct{}

type MatchPlayerData struct {
	ID       int
	JoinTime time.Time
}
type Match_pool struct {
	rwLock         sync.RWMutex
	data           []MatchPlayerData
	MatchTimeRadio float64
}

func (q *Match_pool) GetSize() int {
	q.rwLock.RLock()
	defer q.rwLock.RUnlock()
	return len(q.data)
}

func (q *Match_pool) Add(id int) {
	q.rwLock.Lock()
	defer q.rwLock.Unlock()
	p := MatchPlayerData{ID: id, JoinTime: time.Now()}
	q.data = append(q.data, p)
}

func (q *Match_pool) Delete(id int) {
	q.rwLock.Lock()
	defer q.rwLock.Unlock()
	for i, p := range q.data {
		if p.ID == id {
			q.data[i] = q.data[len(q.data)-1]
			q.data = q.data[:len(q.data)-1]
			break
		}
	}
}

var MatchPool Match_pool

func ImplMatchPool() {
	MatchPool = Match_pool{}
	MatchPool.MatchTimeRadio = global.MatchTimeRadio
	MatchPool.data = make([]MatchPlayerData, 30)
}

//---------------------Match_pool----------------------------

func NewMatchManager() MatchManager {
	ImplMatchPool()
	m := MatchManager{}
	m.StartMatchLoop()
	return m
}

func (m *MatchManager) StartMatchLoop() {
	MatchCheck := time.NewTicker(global.MatchLoopTime * time.Millisecond)
	for {
		select {
		case <-MatchCheck.C:
			m.TryMatch()
		}
	}
}

func (m *MatchManager) AddInPool(id int) {
	MatchPool.Add(id)
}

func (m *MatchManager) TryMatch() bool {

}

func (m *MatchManager) GetRequiredNum() int {
	NowCount := MatchPool.GetSize()
	if NowCount == 0 {
		return global.MatchTimeRadio
	}

}
