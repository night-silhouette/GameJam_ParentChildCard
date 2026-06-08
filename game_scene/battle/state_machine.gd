extends Node
class_name GameManager
@export var card_manager = Node
@export var time = Control
@export var all_block = Control
@export var combat_block = Control
@export var spell_block = Control

@export var judge = Control
@export var jugde_bt = Control
@export var weather = Control
@export var choose_child_card = Control

var is_win:int ;

var parent_action :int ;
var parent_bt_object: int;
var child_action :int ;
var child_bt_object: int;

# 1. 定义状态枚举
enum GameState {
	INIT_STATE,           # 看牌阶段（可拖拽上牌）
	CHOOSE_CHILD_CARD,    # 选择子卡牌阶段
	CHOOSE_WEATHER,       # 选择天气阶段
	USE_MAGIC_CARD,       # 选择技能牌阶段
	USE_COMBAT_CARD,      # 战斗行动阶段
	JUDGEMENT             # 判定阶段（剪刀石头布）
}

# 2. 状态变量（带类型推导与 setter）
var current_state: GameState = GameState.INIT_STATE:
	set(value):
		_exit_state(current_state)
		current_state = value
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


func _enter_state(new_state: GameState) -> void:
	match new_state:
		GameState.INIT_STATE:
			time.countdown_time = TimeOffset.get_remaining_seconds(Global.init_battle_time);
			# print("进入初始化状态（看牌阶段）")
		
		GameState.CHOOSE_CHILD_CARD:
			all_block.block_input()
			Global.revive(choose_child_card)
		
		GameState.CHOOSE_WEATHER:
			all_block.block_input()
			Global.revive(weather)
			weather.init_all_unpressed()
		
		GameState.USE_MAGIC_CARD:
			combat_block.block_input();
			# print("进入魔法卡使用状态")
		
		GameState.USE_COMBAT_CARD:
			spell_block.block_input()
			if(is_win):
				combat_block.block_input();
			SignalBus.request_combat_finish.emit();
			# print("进入战斗卡使用状态")
		
		GameState.JUDGEMENT:
			Global.revive(judge)
			Global.revive(jugde_bt)
			judge.init_all_unpressed()
			combat_block.block_input()
			spell_block.block_input()
			# print("进入结算/判定状态")

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
	time.start_countdown(TimeOffset.get_remaining_seconds(t));
	change_state(GameState.CHOOSE_WEATHER)

# [新增] 天气选择成功
func _on_select_weather_succeed(weather_data) -> void:
	# print("天气选择成功: ", weather_data)
	# 天气选择成功后，进入技能牌选择阶段
	# 实际切换由后端 DeployCard.Query 触发
	pass

# [新增] 卡牌结算完成
func _on_card_calc_finish() -> void:
	# print("卡牌效果结算完成")
	# 结算完成后，进入下一轮的技能牌选择阶段
	# 实际切换由后端信号触发
	pass

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
			SignalBus.request_select_weather.emit(weather.index_judge)
			
		GameState.USE_MAGIC_CARD:
			time.countdown_time = 0;
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.SPELL_ZONE);
			if card.is_empty():
				SignalBus.request_deploy_magic_card.emit(-1,-1);
			else:
				SignalBus.request_deploy_magic_card.emit(card[0].id,card[0].temp_id)
		
		GameState.USE_COMBAT_CARD:
			time.countdown_time = 0;
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.PARENT_BATTLE_ZONE);
			if not card.is_empty():
				SignalBus.request_deploy_parent_card.emit(card[0].id,card[0].temp_id)
			
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.CHILD_BATTLE_ZONE);
			if not card.is_empty():
				SignalBus.request_deploy_child_card.emit(card[0].id,card[0].temp_id)
			
			SignalBus.request_combat_movement.emit(parent_action,0,parent_bt_object);
			SignalBus.request_combat_movement.emit(child_action,0,parent_bt_object);
		
		GameState.JUDGEMENT:
			time.countdown_time = 0;
			jugde_bt.update_single_judge_data(0,judge.index_judge)
			SignalBus.request_judge.emit(judge.index_judge)
			
			

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
