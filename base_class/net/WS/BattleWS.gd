extends Node

var ws := WebSocketPeer.new()
var is_connected := false
func _ready() -> void:
	SignalBus.to_connect_ws.connect(_connect_ws);
func _connect_ws():
	# 假设你的 Token 存在 Global 里
	
	var base_url = Global.BASE_URL.replace("http", "ws") + "/v1/ws/"
	# 👇 将 token 作为 URL 参数拼上去（注意问号）
	var url = base_url + "?token=" + TokenManager.get_token();
	
	var err = ws.connect_to_url(url)
	
	if err != OK:
		push_error("WS连接失败")
	else:
		print("WS连接中... URL: ", url)
		
		
func _process(delta):
	if ws.get_ready_state() == WebSocketPeer.STATE_CONNECTING:
		ws.poll()
	
	elif ws.get_ready_state() == WebSocketPeer.STATE_OPEN:
		if not is_connected:
			is_connected = true
			print("WS连接成功")
			SignalBus.ws_connected.emit()
		
		ws.poll()
		
		while ws.get_available_packet_count() > 0:
			var msg = ws.get_packet().get_string_from_utf8()
			_on_message(msg)
	
	elif ws.get_ready_state() == WebSocketPeer.STATE_CLOSED:
		if is_connected:
			is_connected = false
			print("WS断开")
			SignalBus.ws_disconnected.emit()

# 收消息（只转发，不解析业务）
func _on_message(msg: String):
	var res = JSON.parse_string(msg)
	
	if res == null:
		print("WS解析失败")
		return
	
	var code = res["code"]
	var data = res["data"] if res.has("data") else null
	var msg_str = res["msg"]
	
	# 👇统一出口（和HTTP一样）
	SignalBus.raw_ws_responded.emit(code, data, msg_str)
	
# 发消息（统一入口）
func send_action(action_code: int, action_data = null):
	if ws.get_ready_state() != WebSocketPeer.STATE_OPEN:
		print("WS未连接")
		return
	
	var action = {
		"action_code": action_code,
		"action_name": "",
		"action_data": action_data
	}
	
	var json = JSON.stringify(action)
	ws.send_text(json)
