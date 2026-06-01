extends Control


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	SignalBus.request_bag_card.emit()
	SignalBus.ws_connected.connect(_ws_connected)


# Called every frame. 'delta' is the elapsed time since the previous frame.


func _on_街机按钮_button_down() -> void:
	# 1. 从全局数据层提取需要的原始卡牌数据，过滤掉不需要传给后端的 "zone" 和 "resource"
	var raw_card_list = []
	for card in InventoryManager.get_cards_in_zone(Global.ZONE_CARD.MATCH_ZONE):
		raw_card_list.append({
			"stuff_id": card["stuff_id"],
			"card_id": card["card_id"],
			"price": card["price"]
		})
	var input_text = $"使用/num".text # 获取文本
	var input_amount: int;
	# 1. 判空检查（防止什么都没填就提交）
	if input_text.is_empty():
		#print("请输入数字！")
		input_amount = 0;
	else:
	# 2. 安全转换为 int
		input_amount  = input_text.to_int()
	
	# 3. 验证是否大于 0（防止输入 0 或负数，因为 to_int 转换失败或输入负数会变成≤0）
	if input_amount <0:
		#print("请输入大于 0 的有效数字！")
		return
		
	# 4. 核心：对比全局数据层中的 gold 总数
	# （注意：请把 DataManager 换成你真实的数据层全局单例名）
	var current_gold = InventoryManager.gold 
	
	if input_amount > current_gold:
		#print("大兄弟，金币不足！你当前只有 %d 金币。" % current_gold)
		return
	
	# 2. 严格按照图片中的结构组合成 btData
	var body = {
		"btData": {
			"card_list": raw_card_list,
			"gold": input_amount # 修正拼写错误 glod -> gold，并从数据层获取
		}
	}
	
	# 3. 发送信号并播放等待动画
	SignalBus.to_connect_ws.emit(body)
	$wait.play("wait")

func _on_返回_button_down() -> void:
	SignalBus.change_ui.emit("tomenu")
	
func _ws_connected():
	SignalBus.change_ui.emit("tobattle")
	SignalBus.change_scence.emit("tobattle")
