extends Node
class_name AniStateMachine

# 两个空字典，后续填充动画配置
var ani_enter_dict: Dictionary = {}
var ani_exit_dict: Dictionary = {}
const HP_category = {
	Damage = 0,
	Heal = 1,
	TrueDamage = 2,
}
# 动画状态枚举

enum CalcState {
	IDLE,
	READ,
	SKILL_CARD_NOTIFY,
	WEATHER_NOTIFY,
	BUFF_NOTIFY,
	ACTION_CARD_NOTIFY,
	DEPLOY_CARD_NOTIFY,
	CHILD_BELONG_CHANGE,
	CARD_POS_CHANGE,
	HP_CHANGE,
	BUFF_CHANGE,
	REFRESH_ALL,
	CARD_CALC_FINISH,
}

func _ready() -> void:
	SignalBus.skill_card_notify.connect(_skill_card_notify)
	SignalBus.weather_notify.connect(_weather_notify)
	SignalBus.buff_notify.connect(_buff_notify)
	SignalBus.action_card_notify.connect(_action_card_notify)
	SignalBus.deploy_card_notify.connect(_deploy_card_notify)
	SignalBus.child_belong_change.connect(_child_belong_change)
	SignalBus.card_pos_change.connect(_card_pos_change)
	SignalBus.hp_change.connect(_hp_change)
	SignalBus.buff_change.connect(_buff_change)
	SignalBus.refresh_all.connect(_refresh_all)
	SignalBus.card_calc_finish.connect(_card_calc_finish)
	# 状态变量（带 setter）

var current_state: CalcState = CalcState.IDLE:
	set(value):
		_exit_state(current_state)
		current_state = value
		_enter_state(current_state)
		
func change_state(new_state: CalcState) -> void:
	current_state = new_state
		
func _exit_state(old_state: CalcState) -> void:
	match old_state:
		CalcState.IDLE:
			pass
		CalcState.READ:
			pass
		CalcState.SKILL_CARD_NOTIFY:
			pass
		CalcState.WEATHER_NOTIFY:
			pass
		CalcState.BUFF_NOTIFY:
			pass
		CalcState.ACTION_CARD_NOTIFY:
			pass
		CalcState.DEPLOY_CARD_NOTIFY:
			pass
		CalcState.CHILD_BELONG_CHANGE:
			pass
		CalcState.CARD_POS_CHANGE:
			pass
		CalcState.HP_CHANGE:
			pass
		CalcState.BUFF_CHANGE:
			pass
		CalcState.REFRESH_ALL:
			pass
		CalcState.CARD_CALC_FINISH:
			pass

func _enter_state(new_state: CalcState) -> void:
	match new_state:
		CalcState.IDLE:
			SignalBus.enter_free.emit()
		CalcState.READ:
			var action = Global.cardcalc_animaiton_list.pop_front()
			if action != null:
				action.call()
			else:
				change_state(CalcState.IDLE)
		CalcState.SKILL_CARD_NOTIFY:
			print("进入法术实施")
		CalcState.WEATHER_NOTIFY:
			print("进入天气通知")
		CalcState.BUFF_NOTIFY:
			print("进入Buff通知")
		CalcState.ACTION_CARD_NOTIFY:
			print("进入行动卡通知")
		CalcState.DEPLOY_CARD_NOTIFY:
			print("进入部署卡通知")
		CalcState.CHILD_BELONG_CHANGE:
			print("进入子牌归属变更")
		CalcState.CARD_POS_CHANGE:
			print("进入卡牌位置变更")
		CalcState.HP_CHANGE:
			print("进入HP变更")
		CalcState.BUFF_CHANGE:
			print("进入Buff变更")
		CalcState.REFRESH_ALL:
			print("进入全量刷新")
		CalcState.CARD_CALC_FINISH:
			print("进入卡牌结算完成")

func _skill_card_notify():
	change_state(CalcState.SKILL_CARD_NOTIFY)

func _weather_notify():
	change_state(CalcState.WEATHER_NOTIFY)

func _buff_notify():
	change_state(CalcState.BUFF_NOTIFY)

func _action_card_notify(caller, acceptor, behavior):
	change_state(CalcState.ACTION_CARD_NOTIFY)

func _deploy_card_notify():
	change_state(CalcState.DEPLOY_CARD_NOTIFY)

func _child_belong_change(origin, object):
	change_state(CalcState.CHILD_BELONG_CHANGE)

func _card_pos_change():
	change_state(CalcState.CARD_POS_CHANGE)

func _hp_change():
	change_state(CalcState.HP_CHANGE)

func _buff_change():
	change_state(CalcState.BUFF_CHANGE)

func _refresh_all(All_data):
	change_state(CalcState.REFRESH_ALL)

func _card_calc_finish():
	change_state(CalcState.CARD_CALC_FINISH)
