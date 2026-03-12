package BattleService

import (
	"math"
	"math/rand/v2"
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

func (q *Match_pool) GetMatchTimeRadio() float64 {
	q.rwLock.RLock()
	defer q.rwLock.RUnlock()
	return q.MatchTimeRadio
}
func (q *Match_pool) UpdateMatchTimeRadio(r float64) {
	q.rwLock.Lock()
	defer q.rwLock.Unlock()
	q.MatchTimeRadio = r
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
			q.data = append(q.data[:i], q.data[i+1:]...)
			break
		}
	}
}

func (q *Match_pool) GetMaxHadWaitTime() float64 {
	q.rwLock.RLock()
	defer q.rwLock.RUnlock()
	if len(q.data) == 0 {
		return 0
	}
	return time.Since(q.data[0].JoinTime).Seconds()
}

var MatchPool Match_pool

func ImplMatchPool() {
	MatchPool = Match_pool{}
	MatchPool.MatchTimeRadio = global.MatchTimeRadio
	MatchPool.data = make([]MatchPlayerData, 0, 30)
}

//---------------------Match_pool----------------------------

func NewMatchManager() MatchManager { //对外接口
	ImplMatchPool()
	m := MatchManager{}
	go m.StartMatchLoop()
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

func (m *MatchManager) AddPool(id int) {
	MatchPool.Add(id)
}

func (m *MatchManager) TryMatch() bool {
	NeedNum := m.GetRequiredNum()
	NowNum := MatchPool.GetSize()
	if NowNum >= NeedNum {
		m.MatchCatch()
		return true
	}
	return false
}

func (m *MatchManager) GetRequiredNum() int {
	NowCount := MatchPool.GetSize()
	if NowCount == 0 {
		return 2
	}
	MaxHadWaitTime := MatchPool.GetMaxHadWaitTime()
	m.UpdateMatchTimeRadio(MaxHadWaitTime)
	var res int
	res = int(math.Floor(2.0 * MatchPool.GetMatchTimeRadio()))
	return res
}
func (m *MatchManager) UpdateMatchTimeRadio(MaxHadWaitTime float64) {
	r := global.MatchTimeRadio - (global.MatchTimeRadio-1)*(MaxHadWaitTime/global.MatchMaxWaitTime)
	if r < 1.0 {
		r = 1.0
	}
	MatchPool.UpdateMatchTimeRadio(r)
}

func (m *MatchManager) MatchCatch() (int, int) {
	MatchPool.rwLock.Lock()
	defer MatchPool.rwLock.Unlock()

	poolSize := len(MatchPool.data)
	if poolSize < 2 {
		return -1, -1
	}

	weightList := make([]int, poolSize)
	totalWeight := 0
	for i, data := range MatchPool.data {
		w := int(math.Pow(time.Since(data.JoinTime).Seconds(), 2)) + 1
		weightList[i] = w
		totalWeight += w
	}

	n1 := rand.IntN(totalWeight)
	idx1 := -1
	tempSum1 := 0
	for i, w := range weightList {
		tempSum1 += w
		if tempSum1 > n1 {
			idx1 = i
			break
		}
	}
	newTotalWeight := totalWeight - weightList[idx1]
	n2 := rand.IntN(newTotalWeight)

	idx2 := -1
	tempSum2 := 0
	for i, w := range weightList {
		if i == idx1 {
			continue
		}
		tempSum2 += w
		if tempSum2 > n2 {
			idx2 = i
			break
		}
	}

	id1 := MatchPool.data[idx1].ID
	id2 := MatchPool.data[idx2].ID
	if idx1 < idx2 {
		MatchPool.data = append(MatchPool.data[:idx2], MatchPool.data[idx2+1:]...)
		MatchPool.data = append(MatchPool.data[:idx1], MatchPool.data[idx1+1:]...)
	} else {
		MatchPool.data = append(MatchPool.data[:idx1], MatchPool.data[idx1+1:]...)
		MatchPool.data = append(MatchPool.data[:idx2], MatchPool.data[idx2+1:]...)
	}

	return id1, id2
}
