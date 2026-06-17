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
	DEAD,
	READ,
	
}

# 状态变量（带 setter）
var current_state: CalcState = CalcState.DEAD:
	set(value):
		_exit_state(current_state)
		current_state = value
		_enter_state(current_state)


func _ready() -> void:
	_enter_state(current_state)


func change_state(new_state: CalcState) -> void:
	current_state = new_state


func _exit_state(old_state: CalcState) -> void:
	match old_state:
		AniState.IDLE:
			pass
		AniState.APPEAR:
			pass
		AniState.ATTACK:
			pass
		AniState.SKILL:
			pass
		AniState.DAMAGED:
			pass
		AniState.DEATH:
			pass
		AniState.SWITCH:
			pass
		AniState.DISAPPEAR:
			pass
		AniState.JUDGE_START:
			pass
		AniState.JUDGE_RESULT:
			pass
		AniState.WEATHER:
			pass


func _enter_state(new_state: AniState) -> void:
	match new_state:
		AniState.IDLE:
			pass
		AniState.APPEAR:
			pass
		AniState.ATTACK:
			pass
		AniState.SKILL:
			pass
		AniState.DAMAGED:
			pass
		AniState.DEATH:
			pass
		AniState.SWITCH:
			pass
		AniState.DISAPPEAR:
			pass
		AniState.JUDGE_START:
			pass
		AniState.JUDGE_RESULT:
			pass
		AniState.WEATHER:
			pass
