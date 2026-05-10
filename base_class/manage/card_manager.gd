extends Node
## 导出变量
@export var card_scene: PackedScene = preload("res://base_class/card/card.tscn")
@export var spawn_container: Control
@export_dir var base_path: String = "res://game_data/card/" # 资源存放的基础路径

var card_list :Array = [];#这里的card只掌握数据，不拥有任何的实体
func _ready() -> void:
	SignalBus.self_inhand_updated.connect(_self_inhand_updated)
	SignalBus.bt_oppinfo_updated.connect(_bt_oppinfo_updated)
	SignalBus.bt_selfinfo_updated.connect(_bt_selfinfo_updated)
## 核心功能：通过 ID 生成resoure
func querry_resoure_by_id(card_id: int) -> Resource:
	# 1. 格式化路径，%03d 会将 1 转换为 001，将 12 转换为 012
	var file_name = "card_%03d.tres" % card_id
	var full_path = base_path.path_join(file_name)
	
	# 2. 检查文件是否存在
	if not FileAccess.file_exists(full_path):
		push_error("CardManager: 找不到资源文件 -> " + full_path)
		return null
	
	# 3. 加载资源
	var data = load(full_path) as CardResource
	if data == null:
		push_error("CardManager: 资源加载失败或类型错误 -> " + full_path)
		return null
	
	return data;
	
func _self_inhand_updated(data):
	_update_cards(data,Global.ZONE_CARD.DECK_ZONE)

func _bt_oppinfo_updated(data):
	var bt_skill = [data.get("skill_card_bt")]
	var bt_parent = [data.get("parent_card_bt")]
	var bt_child = [data.get("child_card_bt")]
	_update_cards(bt_skill,Global.ZONE_CARD.ENEMY_SPELL_ZONE)
	_update_cards(bt_parent,Global.ZONE_CARD.ENEMY_PARENT_ZONE)
	_update_cards(bt_child,Global.ZONE_CARD.ENEMY_CHILD_ZONE)

func _bt_selfinfo_updated(data):
	var bt_skill = [data.get("skill_card_bt")]
	var bt_parent = [data.get("parent_card_bt")]
	var bt_child = [data.get("child_card_bt")]
	_update_cards(bt_skill,Global.ZONE_CARD.SPELL_ZONE)
	_update_cards(bt_parent,Global.ZONE_CARD.PARENT_BATTLE_ZONE)
	_update_cards(bt_child,Global.ZONE_CARD.CHILD_BATTLE_ZONE)

func find_card_by_key(cards: Array, target_dict: Dictionary, key_to_match = "temp_id") -> Dictionary:
	if target_dict == null:
		return {}

	var target_value = target_dict.get(key_to_match)

	for card in cards:
		# 防止数组里混入 null
		if card == null:
			continue

		if card.get(key_to_match) == target_value:
			return card
			
	return {}

func _update_cards(data: Array, ZONE):
	# 1. 数据清洗：过滤掉无效的空位（id 为 -1 或数据为空）
	card_list = card_list.filter(func(card): return card != null)
	var valid_new_data = []
	var active_temp_ids = []
	
	for d in data:
		if d != null and d.get("id", -1) != -1:
			valid_new_data.append(d)
			active_temp_ids.append(d.get("temp_id"))

	# 2. 同步数据：处理【添加】和【更新】
	for new_data in valid_new_data:
		_sync_card_data(new_data, ZONE)

	# 3. 区域清理：处理【删除】
	_cleanup_zone(ZONE, active_temp_ids)

## 内部解耦函数：同步单条卡牌数据
func _sync_card_data(new_data: Dictionary, zone):
	var existing_card :Dictionary = find_card_by_key(card_list, new_data, "temp_id")
	
	if !existing_card.is_empty():
		# 更新数据
		existing_card["zone"] = zone
		existing_card["hp"] = new_data.get("hp", 0)
		existing_card["damage"] = new_data.get("damage", 0)
		existing_card["buff_id"] = new_data.get("buff_id", [])
		# 这里可以根据需要扩展更多字段
	else:
		# 只有在 ID 合法时才初始化（双重保障）
		var new_card_dict = _init_single_card(new_data, zone)
		if new_card_dict != null:
			card_list.append(new_card_dict)
	
## 内部解耦函数：清理指定区域中已不存在的卡牌
func _cleanup_zone(zone, active_ids: Array):
	for i in range(card_list.size() - 1, -1, -1):
		var current_card = card_list[i]
		# 关键逻辑：只管辖当前 ZONE 
		if current_card.get("zone") == zone:
			if not current_card.get("temp_id") in active_ids:
				card_list.remove_at(i)		

func _init_single_card(card_data: Dictionary, zone) -> Dictionary:
	
	# 防止空数据
	if card_data == null:
		return {}

	var new_card = card_data.duplicate()

	var resource_data = querry_resoure_by_id(int(card_data.get("id", -1)))

	# 资源不存在
	if resource_data == null:
		push_error("卡牌资源不存在: " + str(card_data))
		return {}

	new_card["is_combat_card"] = resource_data.is_combat_card
	new_card["is_sub_card"] = resource_data.is_sub_card
	new_card["texture"] = resource_data.card_texture
	new_card["zone"] = zone

	return new_card
	
##获取一个区域内的所有卡牌。
func get_cards_by_zone(zone) -> Array:
	var result: Array = []

	for card in card_list:
		# 防止 null
		if card == null:
			continue

		# 区域匹配
		if card.get("zone") == zone:
			result.append(card)

	return result
	
## 修改指定卡牌的区域。
func _change_card_zone(temp_id, new_zone) -> bool:

	for card in card_list:

		# 防止 null
		if card == null:
			continue

		# 找到目标卡
		if card.get("temp_id") == temp_id:

			# 修改 zone
			card["zone"] = new_zone

			return true

	# 没找到
	return false
#游荡对象需要记录原先的zone，但是数据库中的zone应该改变。返回时根据游荡对象的zone去改变
func remove_card_from_view(card_data):
	pass
