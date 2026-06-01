extends Node

var ws := WebSocketPeer.new()
var is_connected := false
func _ready() -> void:
	SignalBus.to_connect_ws.connect(_connect_ws);
func _connect_ws(body):
	
# 1. 先把金币和卡牌的字典转换为 JSON 字符串
	# 这步会把字典变成：{"card_list":[...],"gold":10} 的字符串
	var btData_string = JSON.stringify(body)
	
	# 2. 对 JSON 字符串进行 URL 编码（非常重要！）
	# 因为 JSON 里面有花括号 {}、引号 ""、逗号 ，这些特殊字符在 URL 里是非法的，必须编码
	var encoded_btData = btData_string.uri_encode()
	
	# 3. 获取 Token
	var token_string = TokenManager.get_token()
	
	# 4. 组装基础 URL
	var base_url = Global.BASE_URL.replace("http", "ws") + "/v1/ws/"
	
	# 5. 把 token 和 编码后的 btData 用 ? 和 & 拼接在一起
	var url = base_url + "?token=" + token_string + "&btData=" + encoded_btData
	
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
	print(res)
	if res == null:
		print("WS解析失败")
		return
	
	var code = res["code"]
	var data = res["data"] if res.has("data") else null
	var msg_str = res["msg"]
	
	# 👇统一出口（和HTTP一样）
	SignalBus.raw_ws_responded.emit(code, data, msg_str)
	
# 发消息（统一入口）
func send_action(action_code: int, action_data = null, predicates: int = 2): # 1 对应 Notify 或你的默认值
	if ws.get_ready_state() != WebSocketPeer.STATE_OPEN:
		print("WS未连接，无法发送 Action: ", action_code)
		return
	
	# 构建与后端 Go 结构体完全一致的字典
	var action = {
		"action_code": action_code,
		"action_name": "", # 后端 map 里有，这里通常可以传空或对应的字符串
		"action_data": action_data,
		"predicates": predicates
	}
	
	var json = JSON.stringify(action)
	ws.send_text(json)
	print("send成功")
	# 调试用：打印发送的内容
	# print("向后端发送: ", json)
