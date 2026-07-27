extends Control

var send := 0


func _ready() -> void:
	SignalBus.request_bag_card.emit()
	SignalBus.ws_connected.connect(_ws_connected)
	$"街机按钮".button_down.connect(_on_街机按钮_button_down)
	SignalBus.loot.connect(_on_loot_response)
	
	# 获取战利品一开始假死，数据到后才复活
	Global.fake_death($"获取战利品")
	_start_loot_loop()


func _on_街机按钮_button_down() -> void:
	var raw_card_list = []
	for card in InventoryManager.get_cards_in_zone(Global.ZONE_CARD.MATCH_ZONE):
		raw_card_list.append({
			"stuff_id": int(card["stuff_id"]),
			"card_id": int(card["card_id"]),
			"price": int(card["price"])
		})
	if raw_card_list.size() != 5:
		return
	var input_text = $"使用/num".text
	var input_amount: int
	if input_text.is_empty():
		input_amount = 0
	else:
		input_amount = input_text.to_int()
	if input_amount < 0:
		return
	var current_gold = InventoryManager.gold
	if input_amount > current_gold:
		return
	var body = {
		"btData": {
			"card_list": raw_card_list,
			"gold": input_amount
		}
	}
	if send == 0:
		SignalBus.to_connect_ws.emit(body)
		send = 1
		$wait.visible = true
		$wait.play("wait")
		$ColorRect2.visible = true


func _on_返回_button_down() -> void:
	SignalBus.change_ui.emit("tomenu")

func _ws_connected():
	pass

func _on_万能按钮_button_down() -> void:
	SignalBus.request_cancel_match.emit()
	send = 0
	$wait.visible = false
	$ColorRect2.visible = false


func _on_loot_response(_data) -> void:
	pass


func _start_loot_loop() -> void:
	while true:
		# 1. GET
		SignalBus.request_loot.emit()
		var data = await SignalBus.loot
		if data == null or (data is Array and data.is_empty()):
			print("loot 完成")
			break
		var item: Dictionary = data[0] if data is Array else data
		var loot_id = item.get("loot_id", 0)
		var card_list = item.get("data", [])
		print("card_list: ",card_list)
		InventoryManager.import_loot_card_list(card_list)
		InventoryManager.loot_id = loot_id
		Global.revive($ColorRect3)
		# 数据导入后复活获取战利品
		Global.revive($"获取战利品")
		# 4. 等待 POST 响应
		await SignalBus.login_success
	
	# 循环结束，再次假死
	Global.fake_death($"获取战利品")
	Global.fake_death($ColorRect3)
