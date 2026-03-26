package BattleDto

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func GetActionByWsResByte(byte []byte) Action {
	res := WsResponse{}
	err := json.Unmarshal(byte, &res)
	if err != nil {
		log.Println(err)
	}
	if res.Code == 0 {
		return res.Data
	} else {
		log.Println(res.Code, res.Msg)
		return Action{}
	}
}

type WsResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data Action `json:"data"`
}

type Action struct {
	ActionCode ActionCode `json:"action_code"`
	ActionName string     `json:"action_name"`
	ActionData any        `json:"action_data"`
}

type ActionCode int

const (
	CancelMatch ActionCode = iota
	GetSelfCardInHard
	GetOpponentCardInHard
	OverBattle
	StartBattle
)

var ActionName = map[ActionCode]string{
	CancelMatch:           "取消匹配",
	GetSelfCardInHard:     "获取自己的卡牌信息",
	GetOpponentCardInHard: "获取对手的卡牌信息",
	OverBattle:            "结束战斗",
	StartBattle:           "开始战斗",
}

func NewAction(actionCode ActionCode, ActionData any) Action {
	res := Action{}
	res.ActionCode = actionCode
	res.ActionName = ActionName[actionCode]
	res.ActionData = ActionData
	return res
}

func Send(conn *websocket.Conn, actionCode ActionCode, actionData any) {
	Act := NewAction(actionCode, actionData)
	err := conn.WriteJSON(Act)
	if err != nil {
		log.Println("发送动作失败:", err)
	}
}
