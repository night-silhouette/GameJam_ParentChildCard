extends Node

func _ready():
	SignalBus.raw_ws_responded.connect(_handle_ws_data)

func _handle_ws_data(code: int, data: Variant, msg: String):
	if code != 0:
		# 处理业务级错误
		
		return
	print("msg: ",msg)
	if data == null: return
	if data is String: data = JSON.parse_string(data)

	var action_code = int(data.get("action_code", -1))
	var action_data = data.get("action_data", null)
	var predicate = int(data.get("predicates", 0)) # 获取后端分类

	# 根据分类进行初步过滤或预处理
	match predicate:
		Predicates.NOTIFY:
			print("[WS通知] 收到系统推送到 Action: ", action_code)
		Predicates.RESULT:
			print("[WS结果] action: ", action_code)
		Predicates.QUERY:
			print("[WS查询] action: ",action_code)

	_dispatch(action_code, action_data, predicate)

# 配置表格式：{ action_code: { predicate: [预期类型, 信号] } }
var _dispatch_map = {
	0: { # CancelMatch
		Predicates.RESULT: [TYPE_NIL, SignalBus.match_canceled]
	},
	1: { # GetSelfCardInHard
		Predicates.QUERY: [TYPE_ARRAY, SignalBus.self_cards_updated]
	},
	2: { # GetOpponentCardInHard
		Predicates.QUERY: [TYPE_ARRAY, SignalBus.opponent_cards_updated]
	},
	3: { # OverBattle
		Predicates.NOTIFY: [TYPE_NIL, SignalBus.battle_over]
	},
	4: { # StartBattle
		Predicates.NOTIFY: [TYPE_NIL, SignalBus.battle_started]
	}
}

func _dispatch(action_code: int, action_data: Variant, predicate: int):
	if not _dispatch_map.has(action_code):
		push_warning("未处理的 ActionCode: ", action_code)
		return
		
	var predicate_map = _dispatch_map[action_code]
	
	if not predicate_map.has(predicate):
		push_warning("Action %d 下未处理的 Predicate: %d" % [action_code, predicate])
		return
		
	var config = predicate_map[predicate]
	var expected_type = config[0]
	var target_signal = config[1]

	# 类型校验
	if expected_type != TYPE_NIL and typeof(action_data) != expected_type:
		push_error("数据类型不匹配")
		return

	# 发射信号
	if expected_type == TYPE_NIL:
		target_signal.emit()
	else:
		target_signal.emit(action_data)
		print(action_data);
enum Predicates {
	EMPTY = 0,
	NOTIFY = 1, # 后端主动通知（如：对手出牌了）
	QUERY = 2,  # 客户端查询的返回（如：获取卡组列表）
	RESULT = 3, # 动作执行的结果（如：选牌成功/失败）
	FINISH = 4  # 流程结束
}
