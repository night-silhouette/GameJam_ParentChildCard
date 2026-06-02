extends Node
## 导出变量
@export var card_scene: PackedScene = preload("res://base_class/card/card.tscn")
@export var spawn_container: Control
@export_dir var base_path: String = "res://game_data/card/" # 资源存放的基础路径
@export var state_machine : Node;
signal UI_date_update
signal change_card_zone(temp_id, new_zone)

var free_card_nextzone = null;
var free_card_prevzone = null;
var hover_card :int   =  -1;
var card_list :Array = [];#这里的card只掌握数据，不拥有任何的实体
func _ready() -> void:
	#WS
	SignalBus.self_inhand_updated.connect(_self_inhand_updated)
	SignalBus.bt_oppinfo_updated.connect(_bt_oppinfo_updated)
	SignalBus.bt_selfinfo_updated.connect(_bt_selfinfo_updated)
	change_card_zone.connect(_change_card_zone)
	
	#game
	SignalBus.enter_freecard.connect(_enter_freecard)
	SignalBus.exit_freecard.connect(_exit_freecard)
	SignalBus.detected_area.connect(_detected_area)
	SignalBus.exit_area.connect(_exit_area)
	SignalBus.enter_hover.connect(_enter_hover)
	SignalBus.exit_hover.connect(_exit_hover)
#region 牌的导入

## 核心功能：通过 ID 生成resoure
func querry_resoure_by_id(card_id: int) -> Resource:
	return InventoryManager._find_card_resource_by_id(card_id)
	
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
	
	UI_date_update.emit()
	
## 内部解耦函数：同步单条卡牌数据
func _sync_card_data(new_data: Dictionary, zone):
	# 【核心修复】：查询时，把要对比的 temp_id 也先转成 int，防止 int 和 float 混用导致匹配失败
	var target_temp_id = int(new_data.get("temp_id", -1))
	
	var existing_card : Dictionary = {}
	for card in card_list:
		if card != null and int(card.get("temp_id", -2)) == target_temp_id:
			existing_card = card
			break
	
	if not existing_card.is_empty():
		# 更新数据
		existing_card["zone"] = zone
		existing_card["hp"] = int(new_data.get("hp", 0))
		existing_card["damage"] = int(new_data.get("damage", 0))
		
		var raw_buff = new_data.get("buff_id", [])
		if raw_buff is float or raw_buff is int:
			existing_card["buff_id"] = [int(raw_buff)]
		else:
			existing_card["buff_id"] = Array(raw_buff).map(func(x): return int(x))
	else:
		var new_card_dict = _init_single_card(new_data, zone)
		# 严格校验：确保不是空字典
		if not new_card_dict.is_empty():
			card_list.append(new_card_dict)

## 内部解耦函数：清理指定区域中已不存在的卡牌
func _cleanup_zone(zone, active_ids: Array):
	# 【核心修复】：把活跃的 ID 队列也全转成整型 int 数组，防止 in 操作符因类型不同而误杀
	var clean_active_ids: Array[int] = []
	for id in active_ids:
		clean_active_ids.append(int(id))

	for i in range(card_list.size() - 1, -1, -1):
		var current_card = card_list[i]
		if current_card.get("zone") == zone:
			var current_temp_id = int(current_card.get("temp_id", -1))
			if not current_temp_id in clean_active_ids:
				card_list.remove_at(i)
## 内部解耦函数：初始化单张卡牌，并做严格的类型转换
func _init_single_card(card_data: Dictionary, zone) -> Dictionary:
	if card_data == null or card_data.is_empty():
		return {}

	# 【核心修复】：网络传过来的是浮点数 21.0，必须强制转为整型 int
	var card_id: int = int(card_data.get("id", -1))
	if card_id == -1:
		push_error("CardManager: 收到不合法的卡牌ID (-1)")
		return {}

	# 查询本地资源
	var resource_data = querry_resoure_by_id(card_id)
	if resource_data == null:
		# 这里的报错能帮你精准定位是不是本地漏了哪张牌的配置
		push_error("CardManager: 本地没有对应 [id: " + str(card_id) + "] 的 .tres 资源文件！")
		return {}

	# 深度拷贝网络数据，并规范化数据类型
	var new_card = {}
	new_card["id"] = card_id
	new_card["temp_id"] = int(card_data.get("temp_id", -1))
	new_card["hp"] = int(card_data.get("hp", 0))
	new_card["damage"] = int(card_data.get("damage", 0))
	
	# 处理 buff_id，如果是单值浮点数(0.0)转成数组或规范的整型
	var raw_buff = card_data.get("buff_id", [])
	if raw_buff is float or raw_buff is int:
		new_card["buff_id"] = [int(raw_buff)]
	else:
		new_card["buff_id"] = Array(raw_buff).map(func(x): return int(x))

	# 注入本地配置静态数据
	new_card["resouce"] = resource_data;
	new_card["zone"] = zone

	return new_card
	
#endregion



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

		# 找到目标卡
		if card.get("temp_id") == temp_id:
		
			card["zone"] = new_zone
			
			UI_date_update.emit()
			return true

	# 没找到
	return false
#游荡对象需要记录原先的zone，但是数据库中的zone应该改变。返回时根据游荡对象的zone去改变

func select_card_by_key(value,key_to_match):
	for card in card_list:
		if card[key_to_match] == value:
			return card;
	return {};
		
func _enter_freecard(temp_id,zone):
	free_card_prevzone = zone;
	_change_card_zone(temp_id,Global.ZONE_CARD.FREE_ZONE);
func _exit_freecard(temp_id):
	var icard = select_card_by_key(temp_id,"temp_id")
	if free_card_nextzone == null:
		_change_card_zone(temp_id,free_card_prevzone);
	elif icard["is_combat_card"] == false and free_card_nextzone == Global.ZONE_CARD.SPELL_ZONE:
		match state_machine.current_state:
			state_machine.GameState.USE_MAGIC_CARD:
				for i in get_cards_by_zone(free_card_nextzone):
					i["zone"] = Global.ZONE_CARD.DECK_ZONE;
			
		_change_card_zone(temp_id,free_card_nextzone)
	elif icard["is_combat_card"] == true and icard["is_sub_card"] == false and free_card_nextzone == Global.ZONE_CARD.PARENT_BATTLE_ZONE:
		match state_machine.current_state:
			state_machine.GameState.USE_COMBAT_CARD:
				if state_machine.is_win == 1:
					for i in get_cards_by_zone(free_card_nextzone):
						i["zone"] = Global.ZONE_CARD.DECK_ZONE;
					_change_card_zone(temp_id,free_card_nextzone)
			state_machine.GameState.INIT_STATE:
				for i in get_cards_by_zone(free_card_nextzone):
					i["zone"] = Global.ZONE_CARD.DECK_ZONE;
				_change_card_zone(temp_id,free_card_nextzone)
				
	elif icard["is_combat_card"] == true and icard["is_sub_card"] == true and free_card_nextzone == Global.ZONE_CARD.CHILD_BATTLE_ZONE:
		match state_machine.current_state:
			state_machine.GameState.USE_COMBAT_CARD:
				for i in get_cards_by_zone(free_card_nextzone):
					i["zone"] = Global.ZONE_CARD.DECK_ZONE;
				_change_card_zone(temp_id,free_card_nextzone)
			state_machine.GameState.INIT_STATE:
				for i in get_cards_by_zone(free_card_nextzone):
					i["zone"] = Global.ZONE_CARD.DECK_ZONE;
				_change_card_zone(temp_id,free_card_nextzone)
				
	elif free_card_nextzone == Global.ZONE_CARD.DECK_ZONE :
		_change_card_zone(temp_id,free_card_nextzone)	
	else :
		_change_card_zone(temp_id,free_card_prevzone)
	
func _detected_area(zone):
	free_card_nextzone = zone;
func _exit_area(zone):
	if zone == free_card_nextzone :
		free_card_nextzone = null;

func _enter_hover(temp_id):
	hover_card = temp_id;
	UI_date_update.emit()
func _exit_hover():
	hover_card = -1;
	UI_date_update.emit()
