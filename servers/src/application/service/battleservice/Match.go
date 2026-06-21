package battleservice

import (
	"math"
	"math/rand/v2"
	"pcc_card/Util"
	"pcc_card/application/entity/BattleData"
	"pcc_card/global"
	"pcc_card/infra/repo/userrepo"
	"sync"
	"time"
)

type MatchManager struct {
	repo userrepo.User_repo
}

type MatchPlayerData struct {
	ID       int       `json:"id"`
	JoinTime time.Time `json:"join_time"`
	data     BattleData.EnterBtData
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

func (q *Match_pool) Add(id int, data BattleData.EnterBtData) {
	q.rwLock.Lock()
	defer q.rwLock.Unlock()
	p := MatchPlayerData{ID: id, JoinTime: time.Now(), data: data}
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

func NewMatchManager(repo userrepo.User_repo) { //对外接口
	ImplMatchPool()
	m := MatchManager{}
	m.repo = repo
	go m.StartMatchLoop()
	MM = m
}

func (m *MatchManager) StartMatchLoop() {
	MatchCheck := time.NewTicker(global.MatchLoopTime * time.Millisecond)
	for {
		select {
		case <-MatchCheck.C:
			Ok, id1, id2, CardIdMap, GoldMoreUserId, cList := m.TryMatch()
			if !Ok {
				continue
			} else {

				BTID := BC.AddBattle(id1, id2, CardIdMap, GoldMoreUserId, cList)
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

func (m *MatchManager) AddPool(id int, data BattleData.EnterBtData) {
	MatchPool.Add(id, data)
}

// 这才是匹配算法的集中地,判定是否要把他抓进去
func (m *MatchManager) TryMatch() (bool, int, int, map[int][]int, int, []int) {
	NeedNum := m.GetRequiredNum()
	NowNum := MatchPool.GetSize()
	if NowNum >= NeedNum {
		_, id1, id2, res, GoldMoreUserId := m.MatchCatch() //这一步是在抓了
		cList := RandChildList()

		return true, id1, id2, res, GoldMoreUserId, cList
	}

	return false, -1, -1, nil, -1, make([]int, 0)
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

// MatchCatch 确认了的，抓取的函数
func (m *MatchManager) MatchCatch() (bool, int, int, map[int][]int, int) { //最后一个输出是金币多的那个人的id
	MatchPool.rwLock.Lock()
	defer MatchPool.rwLock.Unlock()

	poolSize := len(MatchPool.data)
	if poolSize < 2 {
		return false, -1, -1, nil, -1
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

	//取确定的人的id
	id1 := MatchPool.data[idx1].ID
	id2 := MatchPool.data[idx2].ID

	//-------返回玩家手牌数据--------

	res := make(map[int][]int)
	for _, e := range MatchPool.data[idx1].data.CardList {
		res[id1] = append(res[id1], e.CardId)
	}
	for _, e := range MatchPool.data[idx2].data.CardList {
		res[id2] = append(res[id2], e.CardId)
	}
	//-------返回玩家手牌数据--------

	//------判别那个人带入的gold更大---
	GoldMoreUserId := id1
	if MatchPool.data[idx1].data.Gold < MatchPool.data[idx2].data.Gold {
		GoldMoreUserId = id2
	}
	//------判别那个人带入的gold更大---

	//------删除带入的卡牌---------
	//for _, dto := range MatchPool.data[idx1].data.CardList {
	//	m.repo.DeleteStuff(context.Background(), m.repo.Get_db(), id1, dto.StuffId)
	//}
	//for _, dto := range MatchPool.data[idx2].data.CardList {
	//	m.repo.DeleteStuff(context.Background(), m.repo.Get_db(), id2, dto.StuffId)
	//}
	//------删除带入的卡牌---------

	//删除带入的金币
	//m.repo.UpdateAssetGold(context.Background(), m.repo.Get_db(),id1,MatchPool.data[idx1].data.Gold)
	//m.repo.UpdateAssetGold(context.Background(), m.repo.Get_db(),id2,MatchPool.data[idx2].data.Gold)

	//删除匹配池子里的元数据
	if idx1 < idx2 {
		MatchPool.data = append(MatchPool.data[:idx2], MatchPool.data[idx2+1:]...)
		MatchPool.data = append(MatchPool.data[:idx1], MatchPool.data[idx1+1:]...)
	} else {
		MatchPool.data = append(MatchPool.data[:idx1], MatchPool.data[idx1+1:]...)
		MatchPool.data = append(MatchPool.data[:idx2], MatchPool.data[idx2+1:]...)
	}

	return true, id1, id2, res, GoldMoreUserId
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

// 随机野生子牌堆
func RandChildList() []int {
	CList := make([]int, 0)
	for i := range global.Lev1Category3Num {
		Id := 2000 + i
		CList = append(CList, Id)
	}
	for i := range global.Lev1Category4Num {
		Id := 3000 + i
		CList = append(CList, Id)
	}
	res := Util.GetRandomElements[int](CList, 10)
	return res
}

//----------------------------------异步通知handler-----------------------------------

var MatchSignals sync.Map

type MatchResult struct {
	BattleID int `json:"battle_id"`
	Opponent int `json:"opponent_id"`
}
