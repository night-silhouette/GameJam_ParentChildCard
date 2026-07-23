extends Node

## --- 调试配置区 ---
@export_group("手动调试信息")
@export var username: String = "test_user_1"
@export var password: String = "123456"
@export var auto_instance_naming: bool = true

## --- 内部状态 ---
var instance_id: int = 0

func _ready() -> void:
	# 1. 禁用 Token 自动保存（确保只在内存）
	TokenManager.enable_save = false
	
	# 2. 物理隔离：自动排列窗口并分配身份
	_setup_multi_instance_layout()
	
	# 3. 信号监听：WS 相关
	SignalBus.ws_connected.connect(_on_ws_connected)
	SignalBus.ws_disconnected.connect(_on_ws_disconnected)
	SignalBus.self_cards_updated.connect(_on_cards_received)
	
	# 4. 信号监听：登录相关（可选，用于打印 UI 反馈）
	# 如果你的 NetworkClient 或 TokenManager 有登录成功的信号，可以在这里连一下
	
	_print_help_menu()

## =========================
## 核心控制逻辑 (手动按键)
## =========================

func _input(event: InputEvent) -> void:
	if event is InputEventKey and event.pressed:
		match event.keycode:
			# --- 第一阶段：HTTP 登录 ---
			KEY_ENTER:
				_step_1_login()
			
			# --- 第二阶段：WS 连接 ---
			KEY_C:
				_step_2_connect_ws()
			
			# --- 第三阶段：WS 业务请求 ---
			KEY_G:
				_step_3_get_cards()
			
			# --- 辅助：手动模拟服务器回包 (不联网时测试分发器) ---
			KEY_M:
				_simulate_server_response()

## =========================
## 业务阶段函数
## =========================

func _step_1_login():
	# print("\n[步骤 1] 发起 HTTP 登录，账号: ", username)
	# 触发你提供的登录逻辑信号
	SignalBus.request_register_user.emit(username, password)
	SignalBus.request_login.emit(username,password)

func _step_2_connect_ws():
	# print("\n[步骤 2] 发起 WS 连接请求 (必须在登录成功后进行)...")
	# 触发连接命令
	SignalBus.to_connect_ws.emit()

func _step_3_get_cards():
	# print("\n[步骤 3] 发送 WS 业务请求：获取手牌...")
	# 触发业务信号
	SignalBus.request_get_self_cards.emit()

## =========================
## WS 回调监听
## =========================

func _on_ws_connected():
	# print("[状态] 成功：WebSocket 已连接！(实例: %d)" % instance_id)
	pass

func _on_ws_disconnected():
	# print("[警告] WebSocket 连接断开")
	pass

func _on_cards_received(cards: Array):
	# print("[最终结果] 业务信号收到数据包内容：", cards)
	pass

## =========================
## 分发器模拟 (模拟后端回包)
## =========================

func _simulate_server_response():
	# print("\n[模拟] 构造服务器回包数据流向 BattleWs...")
	var mock_data = {
		"code": 0,
		"msg": "OK",
		"data": {
			"action_code": 1, # 对应获取手牌
			"action_data": ["雷电术", "治疗术", "护盾"]
		}
	}
	var json_str = JSON.stringify(mock_data)
	
	# 模拟消息到达你的 BattleWs 接收函数
	if has_node("/root/BattleWs"):
		get_node("/root/BattleWs")._on_message(json_str)
	else:
		# 如果是单例
		BattleWs._on_message(json_str)

## =========================
## 物理隔离工具
## =========================

func _setup_multi_instance_layout():
	# 根据窗口初始位置识别实例序号
	var window_pos = DisplayServer.window_get_position()
	instance_id = (window_pos.x / 400) + 1 
	
	if auto_instance_naming:
		username = "test_user_" + str(instance_id)
	
	# 重新排布窗口：长条形，方便并排看控制台打印
	var win_size = Vector2i(450, 700)
	var win_pos = Vector2i((instance_id - 1) * 480 + 50, 100)
	
	DisplayServer.window_set_size(win_size)
	DisplayServer.window_set_position(win_pos)
	DisplayServer.window_set_title("实例 %d | 账号: %s" % [instance_id, username])

func _print_help_menu():
	# print("==========================================")
	# print("   联网多例手动测试工具 (实例: %d)" % instance_id)
	# print("   [Enter] : 执行登录 (获取 Token)")
	# print("   [C]     : 连接 WebSocket")
	# print("   [G]     : 获取手牌 (WS 业务请求)")
	# print("   [M]     : 模拟服务器回包 (本地链路测试)")
	# print("==========================================")
