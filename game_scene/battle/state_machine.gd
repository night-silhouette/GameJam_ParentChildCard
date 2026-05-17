extends Node
class_name GameManager # 或者卡牌管理器等
@export var card_manager = Node
@export var time = Control
@export var all_block = Control
@export var combat_block = Control
@export var spell_block = Control

@export var judge = Control;
@export var jugde_bt = Control;

var is_win:int ;

var parent_action :int ;
var parent_bt_object: int;
var child_action :int ;
var child_bt_object: int;

# 1. 定义状态枚举
enum GameState {
	INIT_STATE,
	USE_MAGIC_CARD,
	USE_COMBAT_CARD,
	JUDGEMENT
}

# 2. 状态变量（带类型推导与 setter）
var current_state: GameState = GameState.INIT_STATE:
	set(value):
		SignalBus.request_get_self_cards_inhand.emit()
		SignalBus.request_get_combat_cards.emit()
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

	# 初始化进入第一个状态
	_enter_state(current_state)
	Global.fake_death(judge)
	Global.fake_death(jugde_bt)
	
	
# --- 核心：状态切换逻辑 ---
func change_state(new_state: GameState) -> void:
	current_state = new_state # 触发 setter

func _exit_state(old_state: GameState) -> void:
	match old_state:
		GameState.INIT_STATE:
			pass
		GameState.USE_MAGIC_CARD:
			combat_block.allow_input()
		
		GameState.USE_COMBAT_CARD:
			spell_block.allow_input()
			
		GameState.JUDGEMENT:
			Global.fake_death(judge)
			Global.fake_death(jugde_bt)
			combat_block.allow_input()
			spell_block.allow_input()


func _enter_state(new_state: GameState) -> void:
	match new_state:
		GameState.INIT_STATE:
			time.countdown_time = TimeOffset.get_remaining_seconds(Global.init_battle_time);
			print("进入初始化状态")
		GameState.USE_MAGIC_CARD:
			combat_block.block_input();
			print("进入魔法卡使用状态")
		GameState.USE_COMBAT_CARD:
			spell_block.block_input()
			SignalBus.request_combat_finish.emit();
			print("进入战斗卡使用状态")
		GameState.JUDGEMENT:
			Global.revive(judge);
			Global.revive(jugde_bt);
			combat_block.block_input();
			spell_block.block_input()
			print("进入结算/判定状态")

# --- 信号回调（在这里控制状态流转） ---
func _on_match_success(t) -> void:
	change_state(GameState.INIT_STATE)
	
func _on_magic_card_start(t) -> void:
	time.countdown_time = TimeOffset.get_remaining_seconds(t);
	change_state(GameState.USE_MAGIC_CARD)
	
func _on_combat_start_success(t,dis_win) -> void:
	time.countdown_time = TimeOffset.get_remaining_seconds(t);
	is_win = int(dis_win);
	change_state(GameState.USE_COMBAT_CARD)
	
func _judge_start(t) -> void:
	time.countdown_time = TimeOffset.get_remaining_seconds(t);
	change_state(GameState.JUDGEMENT) 

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
			SignalBus.request_deploy_parent_card.emit(card[0].id,card[0].temp_id)
			card = card_manager.get_cards_by_zone(Global.ZONE_CARD.CHILD_BATTLE_ZONE);
			SignalBus.request_deploy_child_card.emit(card[0].id,card[0].temp_id)
			SignalBus.request_combat_movement.emit(parent_action,0,parent_bt_object);
			SignalBus.request_combat_movement.emit(child_action,0,child_bt_object);
		GameState.JUDGEMENT:
			time.countdown_time = 0;
			jugde_bt.update_single_judge_data(0,judge.index_judge)
			SignalBus.request_judge.emit(judge.index_judge)
			
			

func _deploy_magic_success():
	if current_state == GameState.USE_MAGIC_CARD :
		pass;
func _magic_card_finish(	):
	if current_state == GameState.USE_MAGIC_CARD :
		pass;
func _judge_finish(data):
	if current_state == GameState.JUDGEMENT:
		jugde_bt.judge_data = [int(data.self), int(data.opponent)]
		jugde_bt.is_win = int(data.is_win)
		
