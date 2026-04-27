extends Node
#
enum State { use_magic_card, judgement , current_round_index , use_combat_card}
var current_state 


func _ready():
	# 所有的状态切换都由信号触发，而不是每帧判断
	SignalBus.match_success.connect(_match_success);
	SignalBus.magic_card_start.connect(_magic_card_start);
	SignalBus.combat_start_success.connect(_combat_start_success)
func _match_success(t):
	pass;
	
func _magic_card_start(t):
	pass;
	
func _combat_start_success(t):
	pass;
	
#func _on_timer_timeout():
	#if current_state == State.IDLE:
		#_transition_to(State.WANDER)
	#else:
		#_transition_to(State.IDLE)
#
#func _on_enemy_spotted(body):
	#if body.is_in_group("player"):
		#_transition_to(State.ALERT)
#
#func _on_enemy_lost(body):
	#if body.is_in_group("player"):
		#_transition_to(State.IDLE)
		#timer.start() # 重新开始巡逻计时
#func _transition_to(new_state):
#
	#if current_state == new_state:
		#return
	#
	#current_state = new_state
	#print("状态切换至: ", State.keys()[new_state])
	#
	## 根据新状态执行一次性的动作
	#match current_state:
		#
		
