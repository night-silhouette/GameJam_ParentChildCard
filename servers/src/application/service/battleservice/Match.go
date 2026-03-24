package battleservice

import (
	"math"
	"math/rand/v2"
	"pcc_card/global"
	"sync"
	"time"
)

type MatchManager struct{}

type MatchPlayerData struct {
	ID       int       `json:"id"`
	JoinTime time.Time `json:"join_time"`
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

func (q *Match_pool) DebugGetMachData() []MatchPlayerData {
	q.rwLock.RLock()
	defer q.rwLock.RUnlock()
	return q.data
}

//---------------------Match_pool----------------------------

var MM MatchManager

func NewMatchManager() { //对外接口
	ImplMatchPool()
	m := MatchManager{}
	go m.StartMatchLoop()
	MM = m
}

func (m *MatchManager) StartMatchLoop() {
	MatchCheck := time.NewTicker(global.MatchLoopTime * time.Millisecond)
	for {
		select {
		case <-MatchCheck.C:
			Ok, id1, id2 := m.TryMatch()
			if !Ok {
				continue
			} else {
				BTID := BC.AddBattle(id1, id2)
				if v, ok := MatchSignals.Load(id1); ok {
					v.(chan MatchResult) <- MatchResult{BattleID: BTID, Opponent: id2}
				}
				if v, ok := MatchSignals.Load(id2); ok {
					v.(chan MatchResult) <- MatchResult{BattleID: BTID, Opponent: id1}
				}
			}

		}
	}
}

func (m *MatchManager) AddPool(id int) {
	MatchPool.Add(id)
}

func (m *MatchManager) TryMatch() (bool, int, int) {
	NeedNum := m.GetRequiredNum()
	NowNum := MatchPool.GetSize()
	if NowNum >= NeedNum {
		_, id1, id2 := m.MatchCatch()
		return true, id1, id2
	}
	return false, -1, -1
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

func (m *MatchManager) MatchCatch() (bool, int, int) {
	MatchPool.rwLock.Lock()
	defer MatchPool.rwLock.Unlock()

	poolSize := len(MatchPool.data)
	if poolSize < 2 {
		return false, -1, -1
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

	return true, id1, id2
}

func (m *MatchManager) IsHasID(id int) bool {
	MatchPool.rwLock.RLock()
	defer MatchPool.rwLock.RUnlock()
	for _, data := range MatchPool.data {
		if data.ID == id {
			return true
		}
	}
	return false
}

//----------------------------------异步通知handler-----------------------------------

var MatchSignals sync.Map

type MatchResult struct {
	BattleID int `json:"battle_id"`
	Opponent int `json:"opponent_id"`
}
