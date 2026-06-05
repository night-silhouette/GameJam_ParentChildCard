extends Node
# 【全局数据层】不挂载任何 UI 节点，只管数据

# 当数据发生变化时，通知全家老小（比如 UI 刷新）
signal bag_updated 
signal gold_updated(value:int)

var card_list: Array = []:
	set(value):
		card_list = value
		# 数据一进来，立刻在局外完成排序
		card_list.sort_custom(func(a, b): return a["price"] > b["price"])
		# 发射信号，告诉 UI 层：“数据更新了，你们该干嘛干嘛”
		bag_updated.emit()

var _gold : int = 0
var chip : int = 0

var gold : int:
	get:
		return _gold
	set(value):
		_gold = value
		gold_updated.emit(_gold)
		

@export_dir var cards_folder_path: String = "res://game_data/card/"


func _ready() -> void:
	# 数据层直接对接网络/系统信号
	SignalBus.get_card_bag.connect(_get_card_bag)
	SignalBus.get_self_gold.connect(_get_self_gold)
	SignalBus.sell_card_success.connect(_sell_card_success)


func _get_card_bag(incoming_card_data: Array):
	var temp_list = []
	for item in incoming_card_data:
		var card_res = _find_card_resource_by_id(int(item["card_id"]))
		var data = {
			"stuff_id": int(item["stuff_id"]),
			"card_id": int(item["card_id"]),
			"price": int(item["price"]),
			"zone": Global.ZONE_CARD.BAG_ZONE,
			"resource": card_res
		}
		temp_list.append(data)
	
	card_list = temp_list


# 提供给外部修改数据的接口
func move_to_sell_zone(stuff_id: int) -> bool:
	for item in card_list:
		if item["stuff_id"] == stuff_id:
			item["zone"] = Global.ZONE_CARD.SELL_ZONE
			# 修改了内部字典，手动触发信号更新 UI
			bag_updated.emit() 
			return true
	return false
func move_to_bag_zone(stuff_id: int) -> bool:
	for item in card_list:
		if item["stuff_id"] == stuff_id:
			item["zone"] = Global.ZONE_CARD.BAG_ZONE
			# 修改了内部字典，手动触发信号更新 UI
			bag_updated.emit() 
			return true
	return false
func move_to_combat_zone(stuff_id: int) -> bool:
	var match_total_count := 0
	var match_sub_card_count := 0
	# 统计当前出战区
	for card in card_list:
		if card["zone"] == Global.ZONE_CARD.MATCH_ZONE:
			match_total_count += 1

			var res: CardResource = card["resource"]
			if res and res.is_sub_card:
				match_sub_card_count += 1
	# 找目标卡
	for card in card_list:
		if card["stuff_id"] != stuff_id:
			continue
		# 已经在出战区
		if card["zone"] == Global.ZONE_CARD.MATCH_ZONE:
			return true
		var res: CardResource = card["resource"]
		# 总数限制
		if match_total_count >= 5:
			# print("出战区已满（最多5张卡牌）")
			return false
		# 子牌限制
		if res and res.is_sub_card and match_sub_card_count >= 2:
			# print("出战区子牌已满（最多2张子牌）")
			return false
		card["zone"] = Global.ZONE_CARD.MATCH_ZONE
		bag_updated.emit()

		return true

	return false

# 静态资源检索器
func _find_card_resource_by_id(target_id: int) -> CardResource:
	var resource_path = cards_folder_path + "/card_" + str(target_id) + ".tres"
	if ResourceLoader.exists(resource_path):
		return load(resource_path) as CardResource
	push_warning("全局数据层：未找到卡牌资源文件: " + resource_path)
	return null
	
func _get_self_gold(data:int):
	gold = data
## 读取卡牌文件夹，通过文件名解析并返回所有 card_id 的数组
func get_all_card_ids_by_filename() -> Array[int]:
	var id_list: Array[int] = []
	
	# 打开目标文件夹
	var dir = DirAccess.open(cards_folder_path)
	if dir:
		dir.list_dir_begin()
		var file_name = dir.get_next()
		
		while file_name != "":
			# 过滤掉文件夹，只处理 .tres 资源文件，且名字匹配前缀
			if !dir.current_is_dir() and file_name.ends_with(".tres") and file_name.begins_with("card_"):
				# 提取数字部分：比如 "card_101.tres" -> "101"
				var id_str = file_name.trim_prefix("card_").trim_suffix(".tres")
				if id_str.is_valid_int():
					id_list.append(id_str.to_int())
			
			file_name = dir.get_next()
		dir.list_dir_end()
	else:
		push_error("全局数据层：无法打开卡牌文件夹路径: " + cards_folder_path)
		
	# 可选：排个序让输出更整齐
	id_list.sort()
	return id_list

func _sell_card_success():
	SignalBus.request_bag_card.emit();
	SignalBus.request_get_self_gold.emit();

func get_cards_in_zone(target_zone: int) -> Array:
	var result = []
	for card in card_list:
		if card["zone"] == target_zone:
			result.append(card)
	return result
	
func get_sell_zone_stuff_ids() -> Array[int]:
	var result: Array[int] = []
	
	for card in card_list:
		if card.get("zone") == Global.ZONE_CARD.SELL_ZONE:
			result.append(int(card["stuff_id"]))
			
	return result
