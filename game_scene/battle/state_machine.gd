extends Node
class_name GameManager

@export var card_manager = Node
@export var time = Control
@export var all_block = Control
@export var combat_block = Control
@export var spell_block = Control

@export var judge = Control
@export var jugde_bt = Control

@export var choose_child_card = GridContainer
@export var weather: Node = null
@export var desk_ui: Node = null

var is_win:int ;

var parent_action :int ;
var parent_bt_object: int;
var child_action :int ;
var child_bt_object: int;
var current_weather:int;

## 子卡选择阶段：服务器返回的可用子卡列表
var _child_card_list: Array = []

## 天气选择阶段：服务器返回的可用天气列表
var _weather_list: Array = []

## 中断选牌阶段：服务器返回的中断数据
var _interrupt_data: Dictionary = {}
var _interrupt_selected: Array = []

# 1. 定义状态枚举
enum GameState {
	INIT_STATE,           # 看牌阶段（可拖拽上牌）
	CHOOSE_CHILD_CARD,    # 选择子卡牌阶段
	CHOOSE_WEATHER,       # 选择天气阶段
	USE_MAGIC_CARD,       # 选择技能牌阶段
	USE_COMBAT_CARD,      # 战斗行动阶段
	JUDGEMENT,            # 判定阶段（剪刀石头布）
	INTERRUPT,            # 中断选牌阶段（死亡/技能触发，上牌消耗能量）
}

# 2. 状态变量（带类型推导与 setter）
var current_state: GameState = GameState.INIT_STATE:
	set(value):
		_exit_state(current_state)
		current_state = value
		print("！！！！！！！！！！！！当前状态为:",current_state)
		_enter_state(current_state)

# --- 初始化 ---
func _ready() -> void:
	# 所有的状态切换都由外部/全局信号触发
	SignalBus.match_success.connect(_on_match_success)
	SignalBus.magic_card_start.connect(_on_magic_card_start)
	SignalBus.combat_start_success.connect(_on_combat_start_success)
	SignalBus.judge_start.connect(_judge_start)
	SignalBus.deploy_magic_success.connect(_deploy_magic_success)
	SignalBus.magic_card_finish.connect(_magic_card_finish)
	SignalBus.judge_finish.connect(_judge_finish)
	
	# [新增] 0.102 新阶段信号
	SignalBus.active_child_card_start.connect(_on_active_child_card_start)
	SignalBus.active_child_card_finish.connect(_on_active_child_card_finish)
	SignalBus.select_weather_start.connect(_on_select_weather_start)
	SignalBus.select_weather_succeed.connect(_on_select_weather_succeed)
	SignalBus.card_calc_finish.connect(_on_card_calc_finish)
	SignalBus.interrupt_start.connect(_on_interrupt_start)
	SignalBus.interrupt_succeed.connect(_on_interrupt_succeed)
	

	# 初始化进入第一个状态
	_enter_state(current_state)
	Global.fake_death(judge)
	Global.fake_death(jugde_bt)
	Global.fake_death(weather)
	Global.fake_death(choose_child_card)
	
	
# --- 核心：状态切换逻辑 ---
func change_state(new_state: GameState) -> void:
	current_state = new_state # 触发 setter

func _exit_state(old_state: GameState) -> void:
	match old_state:
		GameState.INIT_STATE:
			pass
		
		GameState.CHOOSE_CHILD_CARD:
			all_block.allow_input()
			Global.fake_death(choose_child_card)
			card_manager.clear_selection(3)
		
		GameState.CHOOSE_WEATHER:
			all_block.allow_input()
			Global.fake_death(weather)
		
		GameState.USE_MAGIC_CARD:
			combat_block.allow_input()
		
		GameState.USE_COMBAT_CARD:
			spell_block.allow_input()
			combat_block.allow_input()
			
		GameState.JUDGEMENT:
			Global.fake_death(judge)
			Global.fake_death(jugde_bt)
			combat_block.allow_input()
			spell_block.allow_input()
		
		GameState.INTERRUPT:
			all_block.allow_input()
			_interrupt_selected.clear()

func _enter_state(new_state: GameState) -> void:
	# 每进入一个状态，刷新全部战斗数据
	_refresh_all_battle_data()
	
	match new_state:
		GameState.INIT_STATE:
			time.countdown_time = TimeOffset.get_remaining_seconds(Global.init_battle_time);
			all_block.allow_input()
			combat_block.allow_input()
			spell_block.allow_input()
			# print("进入初始化状态（看牌阶段）")
		
		GameState.CHOOSE_CHILD_CARD:
			all_block.block_input()
			Global.revive(choose_child_card)
			_populate_child_cards()
		
		GameState.CHOOSE_WEATHER:
			all_block.block_input()
			Global.revive(weather)
			weather.init_all_unpressed()
			_populate_weather()
		
		GameState.USE_MAGIC_CARD:
			combat_block.block_input();
			# print("进入魔法卡使用状态")
		
		GameState.USE_COMBAT_CARD:
			spell_block.block_input()
			if(is_win):
				combat_block.block_input();
			card_manager.clear_all_combat_dto()
			if desk_ui and desk_ui.has_method("reset_all_combat_cards"):
				desk_ui.reset_all_combat_cards()
			SignalBus.request_combat_finish.emit();
			# print("进入战斗卡使用状态")
		
		GameState.JUDGEMENT:
			Global.revive(judge)
			Global.revive(jugde_bt)
			judge.init_all_unpressed()
			combat_block.block_input()
			spell_block.block_input()
			# print("进入结算/判定状态")
		
		GameState.INTERRUPT:
			all_block.block_input()
			combat_block.block_input()
			spell_block.block_input()
			Global.revive(choose_child_card)
			_populate_interrupt_cards()

func _refresh_all_battle_data() -> void:
	SignalBus.request_get_self_cards_inhand.emit()
	SignalBus.request_get_opponent_cards_inhand.emit()
	SignalBus.request_get_combat_cards.emit()
	SignalBus.request_get_energy.emit()
	SignalBus.request_get_child_card_list.emit()
	SignalBus.request_get_weather.emit()

func _populate_weather() -> void:
	var children = weather.get_children()
	var index = 0
	for child in children:
		if not child is BaseButton:
			continue
		var weather_label = child.get_node_or_null("weather")
		if weather_label == null:
			continue
		if index < _weather_list.size():
			var weather_id = int(_weather_list[index])
			weather_label.text = Global.WEATHER_NAME.get(weather_id, "未知")
			child.visible = true
		else:
			child.visible = false
		index += 1

## 用服务器返回的 child_list 数据填充 choose_child_card 下的每个 choose_card 实例
func _populate_child_cards() -> void:
	var container = choose_child_card
	if not is_instance_valid(container):
		container = get_node_or_null("../../选择子牌")
	if container == null:
		return
	var children = container.get_children()
	for i in range(children.size()):
		var card = children[i]
		if i < _child_card_list.size():
			var data = _child_card_list[i]
			card.match_code = 3  # 子卡激活选择
			card.setup(data)
			card.visible = true
		else:
			card.visible = false

## 用中断数据填充选择子牌界面
func _populate_interrupt_cards() -> void:
	var container = choose_child_card
	if not is_instance_valid(container):
		container = get_node_or_null("../../选择子牌")
	if container == null:
		return
	
	var temp_id_list = _interrupt_data.get("temp_id_list", [])
	if not temp_id_list is Array:
		temp_id_list = []
	var select_num = int(_interrupt_data.get("select_num", 0))
	
	var children = container.get_children()
	# 先把所有卡隐藏
	for child in children:
		child.visible = false
	
	# 填充可选牌：从 hand 中找到匹配 temp_id 的卡牌数据
	var all_cards = card_manager.get_cards_by_zone(Global.ZONE_CARD.HAND_ZONE)
	for i in range(min(temp_id_list.size(), children.size())):
		var tid = int(temp_id_list[i])
		var match_data = {}
		for c in all_cards:
			if int(c.get("temp_id", -1)) == tid:
				match_data = c
				break
		if not match_data.is_empty():
			var card = children[i]
			card.match_code = 99  # 中断选牌
			card.setup(match_data)
			card.visible = true

# --- 信号回调（在这里控制状态流转） ---
func _on_match_success(t) -> void:
	change_state(GameState.INIT_STATE)
	
func _on_magic_card_start(t) -> void:
	time.start_countdown(TimeOffset.get_remaining_seconds(t));
	change_state(GameState.USE_MAGIC_CARD)
	
func _on_combat_start_success(t,dis_win) -> void:
	time.start_countdown(TimeOffset.get_remaining_seconds(t));
	is_win = int(dis_win);
	change_state(GameState.USE_COMBAT_CARD)
	
func _judge_start(t) -> void:
	time.start_countdown(TimeOffset.get_remaining_seconds(t));
	change_state(GameState.JUDGEMENT) 

# [新增] 子卡选择阶段开始
func _on_active_child_card_start(t, child_list) -> void:
	_child_card_list = child_list if child_list is Array else []
	time.start_countdown(TimeOffset.get_remaining_seconds(t));
	change_state(GameState.CHOOSE_CHILD_CARD)

# [新增] 子卡选择阶段结束
func _on_active_child_card_finish(selected_temp_id_list) -> void:
	# 子卡选择结束，进入天气选择阶段（或根据后端流程调整）
	# print("子卡选择结束，选中: ", selected_temp_id_list)
	# 注意：实际状态切换应由后端信号触发，这里只是日志
	pass

# [新增] 天气选择阶段开始
func _on_select_weather_start(t, weather_list) -> void:
	_weather_list = weather_list if weather_list is Array else []
	time.start_countdown(TimeOffset.get_remaining_seconds(t));
	change_state(GameState.CHOOSE_WEATHER)

# [新增] 天气选择成功
func _on_select_weather_succeed(weather_data) -> void:
	pass

# [新增] 卡牌结算完成
func _on_card_calc_finish() -> void:
	pass

# [新增] 中断选牌阶段开始
func _on_interrupt_start(interrupt_data) -> void:
	_interrupt_data = interrupt_data
	_interrupt_selected.clear()
	var t = interrupt_data.get("state_wait_time", 60)
	time.start_countdown(TimeOffset.get_remaining_seconds(t))
	change_state(GameState.INTERRUPT)

# [新增] 中断选牌成功
func _on_interrupt_succeed() -> void:
	_interrupt_selected.clear()

func _on_万能按钮_button_down() -> void:
	var card:Array;
	match current_state :

		GameState.INIT_STATE:
			time.countdown_time = 0;
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.PARENT_BATTLE_ZONE);
			if card.is_empty():
				return
			SignalBus.request_deploy_parent_card.emit(card[0].id,card[0].temp_id)
			
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.SPELL_ZONE);
			if card.is_empty():
				SignalBus.request_deploy_magic_card.emit(-1,-1);
			else:
				SignalBus.request_deploy_magic_card.emit(card[0].id,card[0].temp_id)

		GameState.CHOOSE_CHILD_CARD:
			# 子卡选择阶段：提交选中的子卡
			time.countdown_time = 0;
			var selected = card_manager.get_active_child_list()
			SignalBus.request_active_child_card.emit(selected)
		
		GameState.CHOOSE_WEATHER:
			time.countdown_time = 0
			var weather_id = -1
			if weather.index_judge >= 0 and weather.index_judge < _weather_list.size():
				weather_id = int(_weather_list[weather.index_judge])
			SignalBus.request_select_weather.emit(weather_id)
			
		GameState.USE_MAGIC_CARD:
			time.countdown_time = 0;
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.SPELL_ZONE);
			if card.is_empty():
				SignalBus.request_deploy_magic_card.emit(-1,-1);
			else:
				SignalBus.request_deploy_magic_card.emit(card[0].id,card[0].temp_id)
		
		GameState.USE_COMBAT_CARD:
			time.countdown_time = 0
			var card = card_manager.get_cards_by_zone(Global.ZONE_CARD.PARENT_BATTLE_ZONE)
			if not card.is_empty():
				SignalBus.request_deploy_parent_card.emit(card[0].id, card[0].temp_id)
			
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.CHILD_BATTLE_ZONE)
			if not card.is_empty():
				SignalBus.request_deploy_child_card.emit(card[0].id, card[0].temp_id)
			
			var dp = card_manager.parent_combat_dto
			var dc = card_manager.child_combat_dto
			if dp.behavior != -1:
				SignalBus.request_combat_movement.emit(dp.behavior, dp.self_where, dp.opponent_where, dp.get("select_card", {}))
			if dc.behavior != -1:
				SignalBus.request_combat_movement.emit(dc.behavior, dc.self_where, dc.opponent_where, dc.get("select_card", {}))
		
		GameState.JUDGEMENT:
			time.countdown_time = 0;
			jugde_bt.update_single_judge_data(0,judge.index_judge)
			SignalBus.request_judge.emit(judge.index_judge)
		
		GameState.INTERRUPT:
			time.countdown_time = 0
			var selected = card_manager.get_selected_temp_ids(99)
			if selected.is_empty():
				SignalBus.request_interrupt_select.emit([])
			else:
				# 每选一张牌消耗 1 能量
				var cost = selected.size()
				if card_manager.self_energy >= cost:
					card_manager.self_energy -= cost
					SignalBus.request_interrupt_select.emit(selected)
				else:
					SignalBus.request_interrupt_select.emit([])
			_interrupt_selected.clear()
			
			

func _deploy_magic_success():
	if current_state == GameState.USE_MAGIC_CARD :
		pass;
func _magic_card_finish( ):
	if current_state == GameState.USE_MAGIC_CARD :
		pass;
func _judge_finish(data):
	if current_state == GameState.JUDGEMENT:
		Global.fake_death(judge)
		jugde_bt.judge_data = [int(data.self), int(data.opponent)]
		jugde_bt.is_win = int(data.is_win)
