extends Node

func _ready() -> void:
	print("--- 启动架构全流程测试 ---")
	
	# 1. 监听连接状态信号
	SignalBus.ws_connected.connect(_on_ws_connected)
	SignalBus.ws_disconnected.connect(_on_ws_disconnected)
	
	# 2. 监听业务逻辑信号（最终结果）
	SignalBus.self_cards_updated.connect(_on_cards_received)
	
	# 3. 触发连接命令
	print("[步骤 1] 发起 WS 连接请求...")
	SignalBus.to_connect_ws.emit()

# --- 连接状态测试 ---

func _on_ws_connected():
	print("[步骤 2] 成功：WS 已连接！")
	
	# 连接成功后，立即测试发送请求
	_test_send_request()

func _on_ws_disconnected():
	print("[警告] WS 连接断开或连接失败")

# --- 业务请求测试 ---

func _test_send_request():
	print("[步骤 3] 模拟发送业务请求：获取手牌...")
	# 触发信号映射字典中的信号
	SignalBus.request_get_self_cards.emit()
	
	# 提示：如果你现在没有真实的后端服务器，
	# 你可以在 2 秒后手动模拟一个服务器回包来测试分发器逻辑
	get_tree().create_timer(2.0).timeout.connect(_simulate_server_response)

# --- 分发器逻辑测试（模拟回包） ---

func _simulate_server_response():
	print("[步骤 4] 模拟服务器回包数据流向...")
	
	# 构造符合你协议的消息格式
	var mock_data = {
		"code": 0,
		"msg": "OK",
		"data": {
			"action_code": 1, # 对应你 _dispatch_map 里的获取手牌
			"action_data": ["雷电术", "治疗术", "护盾"]
		}
	}
	
	var json_str = JSON.stringify(mock_data)
	
	# 模拟消息到达 BattleWs 的接收函数
	# 这里假设你的 BattleWs 脚本中有 _on_message 方法
	BattleWs._on_message(json_str)

# --- 最终结果验证 ---

func _on_cards_received(cards: Array):
	print("[步骤 5] 成功：业务信号收到数据包！内容为：", cards)
	print("--- 测试完成：架构链路已跑通 ---")
