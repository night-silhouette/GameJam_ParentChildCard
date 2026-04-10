extends Node

func _ready():
	SignalBus.raw_ws_responded.connect(_handle_ws_data)

func _handle_ws_data(code: int, data: Variant, msg: String):
	if code != 0:
		print("错误码: ", code, NetError.get_message(code))
		return
	
	if data == null:
		return

	# 👇 核心修复逻辑：如果是字符串，就再解析一次
	if data is String:
		data = JSON.parse_string(data)
		
	print("Data 的原始类型是: ", type_string(typeof(data)))	
	
	# 确保现在 data 是字典，避免再次崩溃
	if typeof(data) != TYPE_DICTIONARY:
		print("[错误] Data 格式非法，无法读取 action_code")
		return
	
	var action_code = int(data["action_code"])
	var action_data = data["action_data"]
	
	_dispatch(action_code, action_data)

func _dispatch(action_code: int, action_data):
	# 1. 检查是否存在该协议
	if not _dispatch_map.has(action_code):
		push_warning("未定义的协议: ", action_code)
		return

	var config = _dispatch_map[action_code]
	var expected_type = config[0]
	var target_signal = config[1]

	# 2. 扩展检验逻辑：类型检查
	if expected_type != TYPE_NIL and typeof(action_data) != expected_type:
		push_error("协议 %d 数据类型错误！预期: %d, 实际: %d" % [action_code, expected_type, typeof(action_data)])
		return

	# 3. 发射信号 (根据是否有参数自动适配)
	if expected_type == TYPE_NIL:
		target_signal.emit()
	else:
		target_signal.emit(action_data)

# 定义一个配置表
# 格式: { code: [ 预期类型, 信号 ] }
var _dispatch_map = {
	0: [TYPE_NIL, SignalBus.match_canceled],        # TYPE_NIL 表示不需要参数
	1: [TYPE_ARRAY, SignalBus.self_cards_updated],
	2: [TYPE_ARRAY, SignalBus.opponent_cards_updated],
	3: [TYPE_NIL, SignalBus.battle_over],
	4: [TYPE_NIL, SignalBus.battle_started],
}
