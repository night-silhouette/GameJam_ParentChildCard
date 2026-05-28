extends Node
# 【全局数据层】不挂载任何 UI 节点，只管数据

# 当数据发生变化时，通知全家老小（比如 UI 刷新）
signal bag_updated 

var card_list: Array = []:
	set(value):
		card_list = value
		# 数据一进来，立刻在局外完成排序
		card_list.sort_custom(func(a, b): return a["price"] > b["price"])
		# 发射信号，告诉 UI 层：“数据更新了，你们该干嘛干嘛”
		bag_updated.emit()

@export_dir var cards_folder_path: String = "res://game_data/card/"


func _ready() -> void:
	# 数据层直接对接网络/系统信号
	SignalBus.get_card_bag.connect(_get_card_bag)


func _get_card_bag(incoming_card_data: Array):
	var temp_list = []
	for item in incoming_card_data:
		var card_res = _find_card_resource_by_id(item["card_id"])
		var data = {
			"stuff_id": item["stuff_id"],
			"card_id": item["card_id"],
			"price": item["price"],
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


# 静态资源检索器
func _find_card_resource_by_id(target_id: int) -> CardResource:
	var resource_path = cards_folder_path + "/card_" + str(target_id) + ".tres"
	if ResourceLoader.exists(resource_path):
		return load(resource_path) as CardResource
	push_warning("全局数据层：未找到卡牌资源文件: " + resource_path)
	return null
