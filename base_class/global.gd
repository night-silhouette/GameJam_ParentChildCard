extends Node
var BASE_URL = "http://120.26.145.68:10086" 
# Global.gd 或 NetworkClient.gd
@export var token_save :bool = true;

const ACTION_NAME = {
	1: "取消匹配",
	2: "获取自己的卡牌信息",
	3: "获取对手的卡牌信息",
	4: "获取场上的战斗信息",
	5: "结束战斗",
	6: "开始战斗",
	7: "部署一张牌",
	8: "战斗回合判断",
	9: "匹配成功",
	10: "动画结束",
	11: "执行战斗行动"
}
enum Predicates {
	EMPTY = 0,
	NOTIFY = 1,
	QUERY = 2,
	RESULT = 3,
	FINISH = 4,
	SUCCEED = 5
}
# 配置表格式：{ action_code: { predicate: [预期类型, 信号] } }
const ACTION = {
	CANCEL_MATCH = 1,
	GET_SELF_CARDS = 2,
	GET_OPPONENT_CARDS = 3,
	GET_BT_INFO = 4,
	OVER_BATTLE = 5,
	START_BATTLE = 6,
	DEPLOY_CARD = 7,
	JUDGE = 8,
	MATCH_SUCCESS = 9,
	ANIMATION_END = 10,
	COMBAT = 11,
}
const ZONE_CARD = {
	DECK_ZONE = 1,               # 牌堆
	
	DISCARD_ZONE = 2,       # 弃牌堆
	
	PARENT_BATTLE_ZONE = 3,   # 母牌战斗区
	CHILD_BATTLE_ZONE = 4,    # 子牌战斗区
	SPELL_ZONE = 5,         # 法术牌区
	
	DRAGGING = 6,           # 拖动中
	
	ENEMY_PARENT_ZONE = 7,    # 敌方母牌战斗区
	ENEMY_CHILD_ZONE = 8,     # 敌方子牌区
	ENEMY_SPELL_ZONE = 9    # 敌方法术牌区
}
