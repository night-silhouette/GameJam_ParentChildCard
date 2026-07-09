extends Node
class_name GameManager

@export var card_manager: Node
@export var time: Control
@export var all_block: Control
@export var combat_block: Control
@export var spell_block: Control

@export var judge: Control
@export var jugde_bt: Control

@export var dis_card = Control
@export var choose_child_card: GridContainer
@export var weather: Node = null
@export var desk_ui: Node = null

@export var state_name:Label = null
@export var calc_state_machine:Node = null


var is_win:int ;
var is_pass :int = 0

var parent_action :int ;
var parent_bt_object: int;
var child_action :int ;
var child_bt_object: int;
var current_weather:int;

## 天气选择阶段：服务器返回的可用天气列表
var _weather_list: Array = []
var _interrupt_selected: Array = []

# 1. 定义状态枚举
enum GameState {
	INIT_STATE,           # 看牌阶段（可拖拽上牌）
	CHOOSE_CHILD_CARD,    # 选择子卡牌阶段
	CHOOSE_WEATHER,       # 选择天气阶段
	USE_MAGIC_CARD,       # 选择技能牌阶段
	USE_COMBAT_CARD,      # 战斗行动阶段
	JUDGEMENT,            # 判定阶段（剪刀石头布）
	INTERRUPT,            # 中断选牌阶段（死亡/技能触发，上牌消耗能量）作为尾部，但是自动跳转到CARDCALC，然后进行read
	FREE,
	CARDCALC,			#确定起始和尾部，接收到尾部后进行read
}
const GAME_STATE_NAME = {
	GameState.INIT_STATE:         "看牌阶段",
	GameState.CHOOSE_CHILD_CARD:  "选择子卡",
	GameState.CHOOSE_WEATHER:     "选择天气",
	GameState.USE_MAGIC_CARD:     "法术牌阶段",
	GameState.USE_COMBAT_CARD:    "战斗行动阶段",
	GameState.JUDGEMENT:          "判定阶段",
	GameState.INTERRUPT:          "中断选牌阶段",
	GameState.FREE:               "自由阶段",
	GameState.CARDCALC:           "卡牌效果结算"
}
# 2. 状态变量（带类型推导与 setter）
var current_state: GameState = GameState.INIT_STATE:
	set(value):
		_exit_state(current_state)
		current_state = value
		_enter_state(current_state)
		state_name.text = str(GAME_STATE_NAME[current_state])
		is_pass = 0;
		
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
	SignalBus.enter_free.connect(_enter_free)
	# [新增] 0.102 新阶段信号
	SignalBus.active_child_card_start.connect(_on_active_child_card_start)
	SignalBus.select_weather_start.connect(_on_select_weather_start)
	SignalBus.card_calc_finish.connect(_on_card_calc_finish)
	SignalBus.interrupt_start.connect(_on_interrupt_start)
	SignalBus.interrupt_succeed.connect(_on_interrupt_succeed)	
	Global.fake_death(judge)
	Global.fake_death(jugde_bt)
	Global.fake_death(weather)
	Global.fake_death(choose_child_card)
	_enter_state(current_state)
	_interrupt_selected = card_manager.interrupt_selected
	
# --- 核心：状态切换逻辑 ---
func change_state(new_state: GameState) -> void:
	current_state = new_state # 触发 setter

func _exit_state(old_state: GameState) -> void:
	_send_message()
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
			SignalBus.set_change_lock.emit(false);
			SignalBus.card_use_dead_exit.emit();
			card_manager.clear_cards_need_operate()
			
			
		GameState.CARDCALC:
			SignalBus.request_end_animation.emit()
		GameState.FREE:
			time.visible = true
			pass;

func _enter_state(new_state: GameState) -> void:
	# 每进入一个状态，刷新全部战斗数据
	# ↓↓↓ 临时：测试阶段，所有状态统一 allow_input ↓↓↓
	all_block.allow_input()
	combat_block.allow_input()
	spell_block.allow_input()
	# ↑↑↑ 临时结束 ↑↑↑
	
	match new_state:
		GameState.INIT_STATE:
			time.start_countdown.call_deferred(TimeOffset.get_remaining_seconds(Global.init_battle_time))
			_refresh_all_battle_data()
		
		GameState.CHOOSE_CHILD_CARD:
			_refresh_all_battle_data()
			Global.revive(choose_child_card)
		
		GameState.CHOOSE_WEATHER:
			_refresh_all_battle_data()
			Global.revive(weather)
			weather.init_all_unpressed()
			_populate_weather()
		
		GameState.USE_COMBAT_CARD:
			_refresh_all_battle_data()
			card_manager.clear_all_combat_dto()
			if desk_ui and desk_ui.has_method("reset_all_combat_cards"):
				desk_ui.reset_all_combat_cards()
		GameState.USE_MAGIC_CARD:
			_refresh_all_battle_data()
			
		GameState.JUDGEMENT:
			_refresh_all_battle_data()
			Global.revive(judge)
			judge.init_all_unpressed()
		
		GameState.INTERRUPT:
			_refresh_all_battle_data()
			Global.revive(choose_child_card)
			_interrupt_selected.clear()
			# 匹配可选中卡牌，标记 NEED_OPERATE
			card_manager._on_interrupt_data_received()
			SignalBus.card_use_dead_enter.emit();
			var temp_id_list = card_manager.interrupt_data.get("temp_id_list", [])
			card_manager.mark_cards_need_operate(temp_id_list)
			SignalBus.set_change_lock.emit(true)
			calc_state_machine.change_state(calc_state_machine.CalcState.READ)
		GameState.CARDCALC:
			_refresh_all_battle_data()
			calc_state_machine.change_state(calc_state_machine.CalcState.READ)
		GameState.FREE:
			time.visible = false;

func _refresh_all_battle_data() -> void:
	SignalBus.request_get_self_cards_inhand.emit()
	SignalBus.request_get_opponent_cards_inhand.emit()
	SignalBus.request_get_combat_cards.emit()
	SignalBus.request_get_energy.emit()
	SignalBus.request_get_child_card_list.emit()
	SignalBus.request_get_weather.emit()
	SignalBus.request_get_discard_list.emit()

func _populate_weather() -> void:
	var children = weather.get_children()
	var index = 0
	for child in children:
		if not child is BaseButton:
			continue
		# 找到按钮下的第一个 Label（不管叫什么名字）作为天气名显示
		var weather_label: Label = null
		for c in child.get_children():
			if c is Label:
				weather_label = c
				break
		if weather_label == null:
			continue
		if index < _weather_list.size():
			var weather_id = int(_weather_list[index])
			weather_label.text = Global.WEATHER_NAME.get(weather_id, "未知")
			child.visible = true
		else:
			child.visible = false
		index += 1

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
	
func _on_active_child_card_start(t, child_list) -> void:
	card_manager._on_child_card_list_updated(child_list)
	time.start_countdown(TimeOffset.get_remaining_seconds(t));
	change_state(GameState.CHOOSE_CHILD_CARD)
	
func _on_select_weather_start(t, weather_list) -> void:
	_weather_list = weather_list if weather_list is Array else []

	time.start_countdown(TimeOffset.get_remaining_seconds(t));
	change_state(GameState.CHOOSE_WEATHER)
	
#自由阶段:提前结束的阶段
func _enter_free():
	change_state(GameState.FREE);
# [新增] 子卡选择阶段结束

# [新增] 卡牌结算完成
func _on_card_calc_finish() -> void:
	change_state(GameState.CARDCALC)

# [新增] 中断选牌阶段开始
func _on_interrupt_start(interrupt_data,is_need) -> void:##去控制最后传数组
	card_manager.interrupt_data = interrupt_data
	_interrupt_selected.clear()
	var t = interrupt_data.get("state_wait_time", 60)
	time.start_countdown(TimeOffset.get_remaining_seconds(t))
	change_state(GameState.INTERRUPT)

# [新增] 中断选牌成功
func _on_interrupt_succeed() -> void:
	_interrupt_selected.clear()

func _on_万能按钮_button_down() -> void:
	if is_pass == 1:
		return;
	is_pass = 1;
	SignalBus.enter_free.emit();
		
func _send_message():
	var card:Array;
	match current_state :
		GameState.INIT_STATE:
			
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.PARENT_BATTLE_ZONE);
			if !card.is_empty():
				SignalBus.request_deploy_parent_card.emit(card[0].id,card[0].temp_id)
		
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.CHILD_BATTLE_ZONE)
			if !card.is_empty():
				SignalBus.request_deploy_child_card.emit(card[0].id,card[0].temp_id)
				
		GameState.CHOOSE_CHILD_CARD:
			# 子卡选择阶段：从 CHILD_ACTIVE zone 收集 temp_id
			var selected = card_manager.get_active_child_temp_ids()
			SignalBus.request_active_child_card.emit(selected)
		
		GameState.CHOOSE_WEATHER:
			var weather_id = -1
			if weather.index_judge >= 0 and weather.index_judge < _weather_list.size():
				weather_id = int(_weather_list[weather.index_judge])
			SignalBus.request_select_weather.emit(weather_id)
			
		GameState.USE_MAGIC_CARD:
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.SPELL_ZONE);
			if !card.is_empty():
				SignalBus.request_deploy_magic_card.emit(card[0].id,card[0].temp_id)
			else :
				SignalBus.request_deploy_magic_card.emit(-1,-1)
		
		GameState.USE_COMBAT_CARD:
			var combat_list = card_manager.filter_empty_dicts(card_manager.switch_list + card_manager.action_list)
			print(card_manager.switch_list)
			print(card_manager.action_list)
			print(combat_list)
			SignalBus.request_combat_movement.emit(combat_list)
			
		GameState.JUDGEMENT:
			if judge.index_judge == -1:
				return;
			jugde_bt.update_single_judge_data(0,judge.index_judge)
			SignalBus.request_judge.emit(judge.index_judge)
			
		GameState.INTERRUPT:
			SignalBus.request_interrupt_select.emit(_interrupt_selected);

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
		Global.revive(jugde_bt)
