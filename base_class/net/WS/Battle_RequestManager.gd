extends Node

# 1. 定义信号与 ActionCode 的映射字典
# 格式: { 信号对象: 对应的整数代码 }
@onready var _request_map: Dictionary = {
	SignalBus.request_cancel_match: 0,
	SignalBus.request_get_self_cards: 1,
	SignalBus.request_get_opponent_cards: 2,
	SignalBus.request_over_battle: 3,
}

func _ready():
	# 2. 遍历字典，统一完成连接
	for sig in _request_map:
		var action_code = _request_map[sig]
		# 使用 bind() 将额外的参数(action_code) 绑定到通用处理函数上
		sig.connect(_on_request_received.bind(action_code))

# 3. 统一的处理函数
func _on_request_received(action_code: int):
	BattleWs.send_action(action_code);
