extends Node
## 导出变量
@export var card_scene: PackedScene = preload("res://base_class/card/battle_card.tscn")
@export_dir var base_path: String = "res://game_data/card/" # 资源存放的基础路径
@export var state_machine : Node;
@export var self_energy_node : Control
@export var oppent_energy_node :Control


#region 数据变更信号
signal UI_date_update
signal change_card_zone(temp_id, new_zone)

## [新增] 能量数据变更信号
signal energy_changed(self_energy: int, opponent_energy: int)

## [新增] 天气数据变更信号
signal weather_num_changed(weather_num: int)


## [新增] 弃牌堆数据变更信号
signal discard_changed

## [新增] 战斗行动数据变更信号
signal combat_dto_changed

## [新增] 子牌堆数据准备就绪 → 选择子牌 UI 自动填充
signal child_cards_ready

## [新增] 中断选牌数据准备就绪
signal interrupt_cards_ready(card_list: Array, select_num: int)
#endregion


#region 全局变量（非卡牌数据，用 setter 发信号通知 UI）
## [新增] 己方能量
var self_energy: int = 0:
	set(value):
		self_energy = value
		energy_changed.emit(self_energy, opponent_energy)

## [新增] 对手能量
var opponent_energy: int = 0:
	set(value):
		opponent_energy = value
		energy_changed.emit(self_energy, opponent_energy)

## [新增] 当前天气编号
var weather_num: int = -1:
	set(value):
		weather_num = value
		weather_num_changed.emit(weather_num)

## [新增] 中断选牌数据
var interrupt_data: Dictionary = {};
## 中断选中的 temp_id 列表
var interrupt_selected: Array[int] = []


## 中断选牌：检查是否可以继续添加（由 battle_card 调用）
func can_add_interrupt_selection(temp_id: int) -> bool:
	var limit = int(interrupt_data.get("select_num", 0))
	if limit <= 0:
		return false
	# 已存在则不允许重复添加
	if interrupt_selected.has(temp_id):
		return false
	return interrupt_selected.size() < limit


## 中断选牌：添加 temp_id（调用前需先通过 can_add 校验）
func add_interrupt_selection(temp_id: int) -> void:
	if not interrupt_selected.has(temp_id):
		interrupt_selected.append(temp_id)


## 中断选牌：移除 temp_id
func remove_interrupt_selection(temp_id: int) -> void:
	var idx = interrupt_selected.find(temp_id)
	if idx >= 0:
		interrupt_selected.remove_at(idx)


## 中断选牌：清空选中
func clear_interrupt_selection() -> void:
	interrupt_selected.clear()


## 中断数据到达后：从 hand 匹配 temp_id，发 signal 让 UI 填充
func _on_interrupt_data_received() -> void:
	if interrupt_data.is_empty():
		return
	var temp_id_list = interrupt_data.get("temp_id_list", [])
	if not temp_id_list is Array:
		temp_id_list = []
	var select_num = int(interrupt_data.get("select_num", 0))

	# 从 card_list 中匹配 temp_id，组装可选的卡牌数据
	var matched_cards: Array = []
	for temp_id in temp_id_list:
		for card in card_list:
			if card.get("temp_id") == temp_id:
				matched_cards.append(card.duplicate())
				break
	
	interrupt_cards_ready.emit(matched_cards, select_num)


## [state_machine 调用] 标记指定 temp_id 的卡牌进入 NEED_OPERATE 状态
func mark_cards_need_operate(temp_id_list: Array) -> void:
	for card in card_list:
		var tid = card.get("temp_id", -1)
		card["need_operate"] = tid in temp_id_list


## [state_machine 调用] 清除所有卡牌的 NEED_OPERATE 标记
func clear_cards_need_operate() -> void:
	interrupt_selected.clear()
	for card in card_list:
		card.erase("need_operate")
var combat_list: Array = []; 
#endregion


#region 游戏运行时变量
## 区域重叠栈（解决多 Area 叠加时 free_card_nextzone 被误清的问题）
var _area_stack: Array[int] = []
var free_card_nextzone = null:
	get:
		return _area_stack.back() if not _area_stack.is_empty() else null
var free_card_prevzone = null;
var hover_card :int   =  -1;
var card_list :Array = [];#这里的card只掌握数据，不拥有任何的实体
var action_list:Array = [];
var switch_list:Array = [{},{}];
var switch_child = false;
var switch_parent = false;
var active_child = false;
var active_parent = false;
const behavior = {
	attack = 0,
	skill = 1,
	switch = 2,
}
#endregion


func _ready() -> void:
	# ==================== WS 卡牌数据导入 ====================
	SignalBus.self_inhand_updated.connect(_self_inhand_updated)
	SignalBus.bt_oppinfo_updated.connect(_bt_oppinfo_updated)
	SignalBus.bt_selfinfo_updated.connect(_bt_selfinfo_updated)
	SignalBus.oppent_inhand_updated.connect(_on_oppent_inhand_updated)
	change_card_zone.connect(_change_card_zone)
	SignalBus.discard_list_updated.connect(_on_discard_list_updated)
	SignalBus.child_card_list_updated.connect(_on_child_card_list_updated)
	
	# ==================== [新增] WS 非卡牌数据导入 ====================
	SignalBus.energy_updated.connect(_on_energy_updated)

	SignalBus.weather_update.connect(_on_weather_update)

	
	
	# ==================== 数据变更 → UI ====================
	energy_changed.connect(_on_energy_changed)

	# ==================== 游戏交互信号 ====================
	SignalBus.enter_freecard.connect(_enter_freecard)
	SignalBus.exit_freecard.connect(_exit_freecard)
	SignalBus.detected_area.connect(_detected_area)
	SignalBus.exit_area.connect(_exit_area)
	SignalBus.enter_hover.connect(_enter_hover)
	SignalBus.exit_hover.connect(_exit_hover)


#region ==================== Zone 优先级系统 ====================

## Zone 优先级定义：数值越大优先级越高
## combat zone > hand zone > state zone
const ZONE_PRIORITY = {
	# combat zones (最高优先级)
	Global.ZONE_CARD.PARENT_BATTLE_ZONE: 300,
	Global.ZONE_CARD.CHILD_BATTLE_ZONE: 300,
	Global.ZONE_CARD.SPELL_ZONE: 300,
	Global.ZONE_CARD.ENEMY_PARENT_ZONE: 300,
	Global.ZONE_CARD.ENEMY_CHILD_ZONE: 300,
	Global.ZONE_CARD.ENEMY_SPELL_ZONE: 300,
	
	# hand zones (中等优先级)
	Global.ZONE_CARD.DECK_ZONE: 200,
	Global.ZONE_CARD.ENEMY_HAND_ZONE: 200,
	
	# state zones (最低优先级 - 由 child_state 决定)
	Global.ZONE_CARD.CHILD_ACTIVE: 100,
	Global.ZONE_CARD.CHILD_NOT_ACTIVE: 100,
	Global.ZONE_CARD.CHILD_DIED: 100,
	Global.ZONE_CARD.CHILD_HAS_CATCH: 100,
	
	# 其他
	Global.ZONE_CARD.DISCARD_ZONE: 50,
	Global.ZONE_CARD.FREE_ZONE: 10,
}

## 获取 zone 的优先级，未定义的返回 0
func _get_zone_priority(zone: int) -> int:
	return ZONE_PRIORITY.get(zone, 0)

## 判断新 zone 是否应该覆盖旧 zone
## 规则：新 zone 优先级 >= 旧 zone 优先级时才覆盖
## 例外：FREE_ZONE 总是可以被覆盖
func _should_override_zone(old_zone: int, new_zone: int) -> bool:
	if old_zone == Global.ZONE_CARD.FREE_ZONE:
		return true
	var old_priority = _get_zone_priority(old_zone)
	var new_priority = _get_zone_priority(new_zone)
	return new_priority >= old_priority

#endregion


#region ==================== 卡牌数据导入（原有逻辑） ====================


## 核心功能：通过 ID 查询本地 .tres 资源
func querry_resoure_by_id(card_id: int) -> Resource:
	return InventoryManager._find_card_resource_by_id(card_id)

## 导入：己方手牌数据 → DECK_ZONE (hand zone, 优先级 200)
func _self_inhand_updated(data):
	_update_cards(data, Global.ZONE_CARD.DECK_ZONE)

## 导入：对手战场数据 → 敌方 combat zones (优先级 300)
func _bt_oppinfo_updated(data):
	var bt_skill = [data.get("skill_card_bt")]
	var bt_parent = [data.get("parent_card_bt")]
	var bt_child = [data.get("child_card_bt")]
	_update_cards(bt_skill, Global.ZONE_CARD.ENEMY_SPELL_ZONE)
	_update_cards(bt_parent, Global.ZONE_CARD.ENEMY_PARENT_ZONE)
	_update_cards(bt_child, Global.ZONE_CARD.ENEMY_CHILD_ZONE)

## 导入：己方战场数据 → 己方 combat zones (优先级 300)
func _bt_selfinfo_updated(data):
	var bt_skill = [data.get("skill_card_bt")]
	var bt_parent = [data.get("parent_card_bt")]
	var bt_child = [data.get("child_card_bt")]
	_update_cards(bt_skill, Global.ZONE_CARD.SPELL_ZONE)
	_update_cards(bt_parent, Global.ZONE_CARD.PARENT_BATTLE_ZONE)
	_update_cards(bt_child, Global.ZONE_CARD.CHILD_BATTLE_ZONE)

## 批量同步卡牌数据到指定区域（带优先级判断）
func _update_cards(data: Array, ZONE):
	# 1. 数据清洗：过滤掉无效的空位（id 为 -1 或数据为空）
	card_list = card_list.filter(func(card): return card != null)
	var valid_new_data = []
	var active_temp_ids = []

	for d in data:
		if d != null and d.get("id", -1) != -1:
			valid_new_data.append(d)
			active_temp_ids.append(d.get("temp_id"))

	# 2. 同步数据：处理【添加】和【更新】（带优先级判断）
	for new_data in valid_new_data:
		_sync_card_data(new_data, ZONE)

	# 3. 区域清理：处理【删除】
	_cleanup_zone(ZONE, active_temp_ids)

	UI_date_update.emit()

## 内部解耦函数：同步单条卡牌数据（带优先级判断）
func _sync_card_data(new_data: Dictionary, zone):
	var target_temp_id = int(new_data.get("temp_id", -1))

	var existing_card : Dictionary = {}
	for card in card_list:
		if card != null and int(card.get("temp_id", -2)) == target_temp_id:
			existing_card = card
			break

	if not existing_card.is_empty():
		# 优先级判断：新 zone 优先级 >= 旧 zone 才更新
		var old_zone = int(existing_card.get("zone", 0))
		if not _should_override_zone(old_zone, zone):
			# 不更新 zone，但更新其他数据（hp, damage, buff 等）
			_update_card_stats(existing_card, new_data)
			return
		
		# 更新 zone 和数据
		existing_card["zone"] = zone
		_update_card_stats(existing_card, new_data)
	else:
		var new_card_dict = _init_single_card(new_data, zone)
		if not new_card_dict.is_empty():
			card_list.append(new_card_dict)

## [新增] 仅更新卡牌的数值数据（不更新 zone）
func _update_card_stats(card: Dictionary, new_data: Dictionary):
	card["hp"] = int(new_data.get("hp", 0))
	card["damage"] = int(new_data.get("damage", 0))

	# 支持 buff_list 格式（BuffDto: buff_id, buff_stacks, value）
	var raw_buff = new_data.get("buff_list", new_data.get("buff_id", []))
	if raw_buff is float or raw_buff is int:
		card["buff_list"] = [{"buff_id": int(raw_buff), "buff_stacks": 1, "buff_value": 0.0}]
	elif raw_buff is Array:
		card["buff_list"] = Array(raw_buff).map(func(x):
			if x is Dictionary:
				return {"buff_id": int(x.get("buff_id", 0)), "buff_stacks": int(x.get("buff_stacks", 1)), "buff_value": float(x.get("value", 0.0))}
			else:
				return {"buff_id": int(x), "buff_stacks": 1, "buff_value": 0.0}
		)
	else:
		card["buff_list"] = []
	
	# [新增] 保存 child_state（用于后续优先级判断和 UI 显示）
	if new_data.has("child_state"):
		card["child_state"] = int(new_data.get("child_state", 1))

## 内部解耦函数：清理指定区域中已不存在的卡牌
func _cleanup_zone(zone, active_ids: Array):
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

	var card_id: int = int(card_data.get("id", -1))
	if card_id == -1:
		push_error("CardManager: 收到不合法的卡牌ID (-1)")
		return {}

	var resource_data = querry_resoure_by_id(card_id)
	if resource_data == null:
		push_error("CardManager: 本地没有对应 [id: " + str(card_id) + "] 的 .tres 资源文件！")
		# 不拒绝创建，resouce 留 null，UI 层有 fallback 处理

	var new_card = {}
	new_card["id"] = card_id
	new_card["temp_id"] = int(card_data.get("temp_id", -1))
	new_card["hp"] = int(card_data.get("hp", 0))
	new_card["damage"] = int(card_data.get("damage", 0))
	new_card["form"]= int(card_data.get("form",0))

	# 支持 buff_list 格式（BuffDto: buff_id, buff_stacks, value）
	var raw_buff = card_data.get("buff_list", card_data.get("buff_id", []))
	if raw_buff is float or raw_buff is int:
		new_card["buff_list"] = [{"buff_id": int(raw_buff), "buff_stacks": 1, "buff_value": 0.0}]
	elif raw_buff is Array:
		new_card["buff_list"] = Array(raw_buff).map(func(x):
			if x is Dictionary:
				return {"buff_id": int(x.get("buff_id", 0)), "buff_stacks": int(x.get("buff_stacks", 1)), "buff_value": float(x.get("value", 0.0))}
			else:
				return {"buff_id": int(x), "buff_stacks": 1, "buff_value": 0.0}
		)
	else:
		new_card["buff_list"] = []

	# [新增] 保存 child_state（如果存在）
	if card_data.has("child_state"):
		new_card["child_state"] = int(card_data.get("child_state", 1))

	# 注入本地配置静态数据
	new_card["resouce"] = resource_data
	new_card["zone"] = zone

	return new_card

#endregion


#region ==================== [新增] 非卡牌数据导入 ====================

## [新增] 导入：能量值数据
## 服务器返回 {"self": int, "opponent": int}
func _on_energy_updated(data):
	if data is Dictionary:
		self_energy = int(data.get("self", 0))
		opponent_energy = int(data.get("opponent", 0))

## 能量数据变更后推送到 UI 节点
func _on_energy_changed(_self: int, _opponent: int):
	if is_instance_valid(self_energy_node) and self_energy_node.has_method("_update_lab"):
		self_energy_node.num = _self
	if is_instance_valid(oppent_energy_node) and oppent_energy_node.has_method("_update_lab"):
		oppent_energy_node.num = _opponent

## [新增] 导入：天气数据
func _on_weather_update(weather_num_data):
	weather_num = int(weather_num_data)

## [新增] 导入：子牌堆列表 → 根据 child_state 翻译为对应 zone
## 服务器返回 []ChildCardDto，每张子卡带 child_state 字段
## ChildState: 0=Active, 1=NotActive, 2=Died, 3=HasCatch
## 翻译为 zone：13(CHILD_ACTIVE) ~ 16(CHILD_HAS_CATCH)，child_state + 13
func _on_child_card_list_updated(data):
	if data is Array:
		# 按 child_state 分组，zone = 13 + child_state
		var zone_batches: Dictionary = {}
		for card_data in data:
			var child_state = int(card_data.get("child_state", 1))
			var zone = Global.ZONE_CARD.CHILD_ACTIVE + child_state  # 13 + 0~3
			if not zone_batches.has(zone):
				zone_batches[zone] = []
			zone_batches[zone].append(card_data)
		
		for zone in zone_batches:
			_update_cards(zone_batches[zone], zone)
		
		# 数据更新完毕，通知选择子牌 UI 填充
		child_cards_ready.emit()

## [新增] 导入：弃牌堆数据 → DISCARD_ZONE（显示区域，优先级 50）
## 服务器返回 []CardDto，zone 只管显示
func _on_discard_list_updated(data):
	if data is Array:
		_update_cards(data, Global.ZONE_CARD.DISCARD_ZONE)
		discard_changed.emit()

## 导入：敌方手牌 → ENEMY_HAND_ZONE（zone 17）
func _on_oppent_inhand_updated(data):
	if data is Array:
		_update_cards(data, Global.ZONE_CARD.ENEMY_HAND_ZONE)
		
func load_all_data(data_all):
	var bt_card_info = data_all.get("bt_card_info")
	var self_data = bt_card_info.get("self");
	var opp_data = bt_card_info.get("opponent");
	_bt_selfinfo_updated(self_data);
	_bt_oppinfo_updated(opp_data);
	
	var card_in_hand = data_all.get("card_in_hand")
	_self_inhand_updated(card_in_hand)
	
	var energy = data_all.get("energy")
	_on_energy_updated(energy)
	
	var discard_pool = data_all.get("discard_pool")
	_on_discard_list_updated(discard_pool)
	
	var child_child_pool = data_all.get("child_child_pool")
	_on_child_card_list_updated(child_child_pool)
	
#endregion


#region ==================== 选中触发逻辑（choose_card 类卡牌使用） ====================


signal selection_changed(match_code: int, temp_id: int, is_selected: bool)

## [新增] 维护各匹配码下的选中卡牌集合
## key: match_code(int), value: Array[int] (temp_id 列表)
var selection_pools: Dictionary = {}

## 子卡激活选择上限（match_code=3）
const MAX_CHILD_SELECT: int = 5

var active_card_list#上限为五张，右键取消左键确认，

func toggle_selection(match_code: int, temp_id: int) -> bool:
	if not selection_pools.has(match_code):
		selection_pools[match_code] = []
	
	var pool: Array = selection_pools[match_code]
	var index = pool.find(temp_id)
	
	if index >= 0:
		# 已选中 → 取消选中，zone 恢复为 CHILD_NOT_ACTIVE
		pool.remove_at(index)
		if match_code == 3:
			_change_card_zone(temp_id, Global.ZONE_CARD.CHILD_NOT_ACTIVE)
		selection_changed.emit(match_code, temp_id, false)
		return false
	else:
		# 子卡激活选择有上限
		if match_code == 3 and pool.size() >= MAX_CHILD_SELECT:
			return false
		# 未选中 → 加入选中，zone 改为 CHILD_ACTIVE
		pool.append(temp_id)
		if match_code == 3:
			_change_card_zone(temp_id, Global.ZONE_CARD.CHILD_ACTIVE)
		selection_changed.emit(match_code, temp_id, true)
		return true

## [新增] 获取指定 match_code 下的所有选中 temp_id
func get_selected_temp_ids(match_code: int) -> Array[int]:
	if not selection_pools.has(match_code):
		return []
	return selection_pools[match_code].duplicate()

## [新增] 获取子卡激活选择（match_code=3）的选中列表，供提交用
func get_active_child_list() -> Array:
	return get_selected_temp_ids(3)

## [新增] 从 CHILD_ACTIVE zone 直接获取 temp_id 列表（万能按钮提交用）
func get_active_child_temp_ids() -> Array:
	var result: Array = []
	for card in card_list:
		if card != null and card.get("zone") == Global.ZONE_CARD.CHILD_ACTIVE:
			result.append(card.get("temp_id"))
	return result

## [新增] 清空指定 match_code 的选中状态
func clear_selection(match_code: int) -> void:
	if selection_pools.has(match_code):
		var pool: Array = selection_pools[match_code]
		for temp_id in pool:
			selection_changed.emit(match_code, temp_id, false)
		selection_pools[match_code] = []

## [新增] 清空所有选中状态
func clear_all_selections():
	for match_code in selection_pools.keys():
		clear_selection(match_code)
		
func set_cards_need_operate(temp_ids: Array) -> void:
	# 再把传入数组中的 temp_id 对应的卡牌设为 NEED_OPERATE
	for tid in temp_ids:
		var card = select_card_by_key(int(tid), "temp_id")
		#这里是数据层，有问题，不能这样写，统一在ui_update里
		

#endregion


#region ==================== 区域查询工具 ====================

## 获取一个区域内的所有卡牌
func get_cards_by_zone(zone) -> Array:
	var result: Array = []
	for card in card_list:
		if card == null:
			continue
		if card.get("zone") == zone:
			result.append(card)
	return result

## 修改指定卡牌的区域
func _change_card_zone(temp_id, new_zone) -> bool:
	for card in card_list:
		if card.get("temp_id") == temp_id:
			card["zone"] = new_zone
			UI_date_update.emit()
			return true
	return false

## 通过 key 查找卡牌
func select_card_by_key(value, key_to_match):
	for card in card_list:
		if card[key_to_match] == value:
			return card
	return {}

## 过滤掉数组中的空字典
func filter_empty_dicts(arr: Array) -> Array:
	var result: Array = []
	for item in arr:
		if item is Dictionary and not item.is_empty():
			result.append(item)
	return result

#endregion


#region ==================== 战斗行动管理（CombatData） ====================

## 母牌战斗 DTO（behavior=-1 表示无操作）
var parent_combat_dto: Dictionary = {"behavior": -1, "self_where": 0, "opponent_where": -1, "temp_id": -1, "select_card": {}}

## 子牌战斗 DTO
var child_combat_dto: Dictionary = {"behavior": -1, "self_where": 1, "opponent_where": -1, "temp_id": -1, "select_card": {}}

var swtich_combat_dto:Dictionary = {"behavior": -1, "self_where": -1, "opponent_where": -1, "temp_id": -1, "select_card": {}}

const COST_SWITCH: int = 1
const COST_ATTACK_OR_SKILL: int = 2


## 设置某个 DTO 的操作（由 桌面.gd / _deploy_card_to_zone 调用）
## select_card: 仅换牌(behavior=2)时需要，格式 {where, card_id, card_temp_id}
## 返回 true 表示设置成功（能量够），false 表示能量不足
func set_combat_dto(self_where: int, behavior: int, opponent_where: int, temp_id: int, select_card: Dictionary = {}) -> bool:
	var cost = COST_SWITCH if behavior == 2 else COST_ATTACK_OR_SKILL
	if self_energy < cost:
		return false
	
	self_energy -= cost
	
	if self_where == 0:
		parent_combat_dto = {"behavior": behavior, "self_where": 0, "opponent_where": opponent_where, "temp_id": temp_id, "select_card": select_card}
		action_list.append(parent_combat_dto)
		active_parent = true
	else:
		child_combat_dto = {"behavior": behavior, "self_where": 1, "opponent_where": opponent_where, "temp_id": temp_id, "select_card": select_card}
		action_list.append(child_combat_dto)
		active_child = true
	combat_dto_changed.emit()
	return true



## 清空全部 DTO 并返还全部能量（桌面万能按钮）
func clear_all_combat_dto() -> void:
	action_list = [];
	switch_list = [{},{}];
	switch_child = false;
	switch_parent = false;
	active_child = false;
	active_parent = false;
	SignalBus.request_get_combat_cards.emit();
	SignalBus.request_get_energy.emit();
	SignalBus.request_get_self_cards_inhand.emit()
	
#endregion


#region ==================== 游戏交互逻辑 ====================

func _enter_freecard(temp_id, zone):
	_area_stack.clear()
	free_card_prevzone = zone;
	_change_card_zone(temp_id, Global.ZONE_CARD.FREE_ZONE);

func _exit_freecard(temp_id):
	var icard = select_card_by_key(temp_id, "temp_id")
	if icard.is_empty():
		_change_card_zone(temp_id, free_card_prevzone)
		_area_stack.clear()
		return

	if free_card_nextzone == null:
		_change_card_zone(temp_id, free_card_prevzone)
		_area_stack.clear()
		return

	var can_deploy = _can_deploy_card(icard, free_card_nextzone)
	if can_deploy:
		_deploy_card_to_zone(temp_id, free_card_nextzone)
	else:
		_change_card_zone(temp_id, free_card_prevzone)
	_area_stack.clear()

func _can_deploy_card(icard: Dictionary, target_zone: int) -> bool:
	if not state_machine:
		return false
	
	var res = icard.get("resouce")
	if res == null:
		return false
	
	match state_machine.current_state:
		state_machine.GameState.INIT_STATE:
			if target_zone == Global.ZONE_CARD.DECK_ZONE:
				return true
			if res.is_combat_card :
				if target_zone == Global.ZONE_CARD.CHILD_BATTLE_ZONE and res.is_sub_card:
					return true
				elif target_zone == Global.ZONE_CARD.PARENT_BATTLE_ZONE and !res.is_sub_card:
					return true
			else:		
				return false
	
		state_machine.GameState.USE_MAGIC_CARD:
			if target_zone == Global.ZONE_CARD.DECK_ZONE:
				return true
			if not res.is_combat_card and target_zone == Global.ZONE_CARD.SPELL_ZONE:
				return true
			return false
		
		state_machine.GameState.USE_COMBAT_CARD:
			if res.is_combat_card:
				if target_zone == Global.ZONE_CARD.DECK_ZONE:
					if res.is_sub_card:
						if active_child:
							return false
						if  switch_child :
							return true
						if self_energy >= 1 :
							self_energy -= 1
							switch_child = true
							return true
					else:
						if active_parent:
							return false
						if  switch_parent :
							return true
						if self_energy >= 1 :
							self_energy -= 1
							switch_parent = true
							return true
				if target_zone == Global.ZONE_CARD.CHILD_BATTLE_ZONE and res.is_sub_card :
					if active_child:
						return false
					if  switch_child :
						return true
					if self_energy >= 1 :
						self_energy -= 1
						switch_child = true
						return true
				elif target_zone == Global.ZONE_CARD.PARENT_BATTLE_ZONE and !res.is_sub_card:
					if active_parent:
						return false
					if  switch_parent :
						return true
					if self_energy >= 1 :
						self_energy -= 1
						switch_parent = true
						return true
			return false
			
		state_machine.GameState.INTERRUPT:
			if res.is_combat_card and !get_cards_by_zone(target_zone).is_empty():
				if target_zone == Global.ZONE_CARD.CHILD_BATTLE_ZONE and res.is_sub_card:
					return true
				elif target_zone == Global.ZONE_CARD.PARENT_BATTLE_ZONE and !res.is_sub_card:
					return true
				
	return false

func _deploy_card_to_zone(temp_id: int, target_zone: int):
	# 先把目标位置原有的牌送回牌库
	for card in get_cards_by_zone(target_zone):
		card["zone"] = Global.ZONE_CARD.DECK_ZONE
	_change_card_zone(temp_id, target_zone)
	
	# 战斗阶段手牌→战斗区换牌：自动记录 Switch DTO
	if state_machine and state_machine.current_state == state_machine.GameState.USE_COMBAT_CARD:
		var icard = select_card_by_key(temp_id,"temp_id")
		var res = icard.get("resouce")
		var select_card:Dictionary = {
			"where": -1,
			"card_id" : -1 ,
			"card_temp_id" : -1,
		}
		
		if res.is_sub_card:
			if target_zone == Global.ZONE_CARD.DECK_ZONE:
				select_card = {"where": Global.WHERE.InHand,"card_id" : int(res.id) ,"card_temp_id" : int(temp_id)}
			elif target_zone == Global.ZONE_CARD.CHILD_BATTLE_ZONE:
				select_card = {"where": Global.WHERE.ChildCard,"card_id" : int(res.id) ,"card_temp_id" : int(temp_id)}
			swtich_combat_dto = {"behavior": behavior.switch, "self_where": -1, "opponent_where": -1, "temp_id": -1, "select_card": select_card}
			switch_list[1] = swtich_combat_dto;
		else:
			if target_zone == Global.ZONE_CARD.DECK_ZONE:
				select_card = {"where": Global.WHERE.InHand,"card_id" : int(res.id) ,"card_temp_id" : int(temp_id)}
			elif target_zone == Global.ZONE_CARD.PARENT_BATTLE_ZONE:
				select_card = {"where": Global.WHERE.ParentCard,"card_id" : int(res.id) ,"card_temp_id" : int(temp_id)}
				
			swtich_combat_dto = {"behavior": behavior.switch, "self_where": -1, "opponent_where": -1, "temp_id": -1, "select_card": select_card}
			switch_list[0] = swtich_combat_dto;
			print(switch_list[0])
			
func _detected_area(zone):
	_area_stack.append(zone)

func _exit_area(zone):
	var idx = _area_stack.find(zone)
	if idx >= 0:
		_area_stack.remove_at(idx)

func _enter_hover(temp_id):
	if hover_card == temp_id:
		return
	hover_card = temp_id;
	UI_date_update.emit()

func _exit_hover():
	if hover_card == -1:
		return
	hover_card = -1;
	UI_date_update.emit()

#endregion
