extends Node
class_name AniStateMachine
@export var card_manager : Node
var cardcalc_data;
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
	SignalBus.ani_end.connect(_ani_end)
	# 状态变量（带 setter）

var current_state: CalcState = CalcState.IDLE:
	set(value):
		_exit_state(current_state)
		current_state = value
		_enter_state(current_state)
		
func change_state(new_state: CalcState,action_data = null) -> void:
	cardcalc_data = action_data
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

func _enter_state(new_state: CalcState) -> void:
	print("cardcalc")
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
			SignalBus.ani_skill_card_notify_enter.emit()

		CalcState.WEATHER_NOTIFY:
			print("进入天气通知")
			SignalBus.ani_weather_notify_enter.emit()

		CalcState.BUFF_NOTIFY:
			print("进入Buff通知")
			SignalBus.ani_buff_notify_enter.emit()
		CalcState.ACTION_CARD_NOTIFY:
			print("进入行动卡通知")
			var caller = cardcalc_data.get("caller")
			var acceptor = cardcalc_data.get("acceptor") #都是tempid
			var behavior = cardcalc_data.get("animation_behavior")
			SignalBus.ani_action_card_notify_enter.emit(caller, acceptor, behavior)

		CalcState.DEPLOY_CARD_NOTIFY:
			print("进入部署卡通知")
			SignalBus.ani_deploy_card_notify_enter.emit(cardcalc_data)

		CalcState.CHILD_BELONG_CHANGE:
			print("进入子牌归属变更")
			var origin = cardcalc_data.get("origin") ##来源，三种来源，子牌堆，我方手牌，敌方手牌的枚举
			var object = cardcalc_data.get("object")
			SignalBus.ani_child_belong_change_enter.emit(origin, object)
		CalcState.CARD_POS_CHANGE:
			print("进入卡牌位置变更")
			var object = cardcalc_data.get("object") #where
			var temp_id = int(cardcalc_data.get("temp_id"))
			SignalBus.ani_card_pos_change_enter.emit(object, temp_id)

		CalcState.HP_CHANGE:
			var temp_id = int(cardcalc_data.get("temp_id"))
			var category = int(cardcalc_data.get("category"))#HP_category
			var value = int(cardcalc_data.get("value"))
			print("进入HP变更")
			SignalBus.ani_hp_change_enter.emit(temp_id, category, value)

		CalcState.BUFF_CHANGE:
			print("进入Buff变更")
			SignalBus.ani_buff_change_enter.emit()

		CalcState.REFRESH_ALL:
			card_manager.load_all_data(cardcalc_data)
			print("进入全量刷新")
			SignalBus.ani_end.emit.call_deferred()

func _skill_card_notify():
	change_state(CalcState.SKILL_CARD_NOTIFY)

func _weather_notify():
	change_state(CalcState.WEATHER_NOTIFY)

func _buff_notify():
	change_state(CalcState.BUFF_NOTIFY)

func _action_card_notify(action_data):
	change_state(CalcState.ACTION_CARD_NOTIFY,action_data)

func _deploy_card_notify():
	change_state(CalcState.DEPLOY_CARD_NOTIFY)

func _child_belong_change(action_data):
	change_state(CalcState.CHILD_BELONG_CHANGE,action_data)

func _card_pos_change(action_data):
	change_state(CalcState.CARD_POS_CHANGE,action_data)

func _hp_change(action_data):
	change_state(CalcState.HP_CHANGE,action_data)

func _buff_change():
	change_state(CalcState.BUFF_CHANGE)

func _refresh_all(All_data):
	change_state(CalcState.REFRESH_ALL,All_data)

func _ani_end():
	change_state.call_deferred(CalcState.READ)
