package BattleDto

type Action struct {
	ActionCode ActionCode `json:"action_code"`
	ActionName string     `json:"action_name"`
	ActionData any        `json:"action_data"`
	Predicates Predicates `json:"predicates"`
}

type ActionCode int

const (
	CancelMatch ActionCode = iota
	GetSelfCardInHard
	GetOpponentCardInHard
	OverBattle
	StartBattle
	SelectCharacterCard
)

var ActionName = map[ActionCode]string{
	CancelMatch:           "取消匹配",
	GetSelfCardInHard:     "获取自己的卡牌信息",
	GetOpponentCardInHard: "获取对手的卡牌信息",
	OverBattle:            "结束战斗",
	StartBattle:           "开始战斗",
	SelectCharacterCard:   "选择角色牌",
}

type Predicates int

const (
	Empty Predicates = iota
	Notify
	Query
	Result
	finish
)

func NewAction(actionCode ActionCode, Predicates Predicates, ActionData any) Action {
	res := Action{}
	res.ActionCode = actionCode
	res.ActionName = ActionName[actionCode]
	res.ActionData = ActionData
	res.Predicates = Predicates
	return res
}
