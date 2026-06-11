extends "res://base_class/card/battle_card.gd"

#region 战斗状态枚举
enum CombatState {
	IDLE,       # 默认状态
	SELECTED,   # 己方牌被选中，等待选目标
	ACTED,      # 本回合已行动，不可再点击
}
enum CombatAnim {
	NONE,
	ATTACKING,
	SKILL,
	DAMAGED,
	DEATH,
}
#endregion

var combat_state: CombatState = CombatState.IDLE
var combat_anim: CombatAnim = CombatAnim.NONE

## 自己的位置（Where）：0=母牌, 1=子牌
var self_where: int = -1

## 是否为己方牌（true=己方可操作, false=敌方仅可被选为目标）
var is_own: bool = false

#region 信号
signal combat_selected(temp_id: int, self_where: int)
signal combat_target_changed(opponent_where: int)
signal switch_card_requested(temp_id: int, self_where: int)
#endregion

#region 视觉节点
@onready var highlight: ColorRect = $"Highlight" if has_node("Highlight") else null
@onready var state_label: Label = $"StateLabel" if has_node("StateLabel") else null
#endregion


func _ready() -> void:
	super._ready()
	
	# 根据 zone 判断是否为己方牌
	if zone in [Global.ZONE_CARD.PARENT_BATTLE_ZONE, Global.ZONE_CARD.CHILD_BATTLE_ZONE]:
		is_own = true
		self_where = 0 if zone == Global.ZONE_CARD.PARENT_BATTLE_ZONE else 1


## 覆盖父类 update_card_data，补充 combat 相关初始化
func update_card_data(base_res: Dictionary) -> void:
	super.update_card_data(base_res)
	
	zone = base_res.get("zone", zone)
	if zone in [Global.ZONE_CARD.PARENT_BATTLE_ZONE, Global.ZONE_CARD.CHILD_BATTLE_ZONE]:
		is_own = true
		self_where = 0 if zone == Global.ZONE_CARD.PARENT_BATTLE_ZONE else 1
	else:
		is_own = false
		self_where = -1


## [核心] 重置为未行动状态（每回合开始调用）
func reset_combat() -> void:
	combat_state = CombatState.IDLE
	combat_anim = CombatAnim.NONE
	_update_visual()


## --- 鼠标交互（覆盖父类） ---
func _gui_input(event: InputEvent) -> void:
	if combat_state == CombatState.ACTED:
		return
	
	if event is InputEventMouseButton and event.pressed and event.button_index == MOUSE_BUTTON_LEFT:
		accept_event()
		if is_own:
			_on_own_left_click()
		else:
			_on_enemy_click()


## 左键点击己方牌 → 切换选中/取消选中
func _on_own_left_click() -> void:
	match combat_state:
		CombatState.IDLE:
			combat_state = CombatState.SELECTED
			combat_selected.emit(temp_id, self_where)
			_update_visual()
		
		CombatState.SELECTED:
			combat_state = CombatState.IDLE
			combat_selected.emit(-1, -1)
			_update_visual()


## 左键点击敌方牌 → 确认目标
func _on_enemy_click() -> void:
	var opponent_where = 0 if zone == Global.ZONE_CARD.ENEMY_PARENT_ZONE else 1
	combat_target_changed.emit(opponent_where)


## 确认行动后变为 ACTED
func mark_acted() -> void:
	combat_state = CombatState.ACTED
	_update_visual()


## 播放战斗动画
func play_combat_anim(anim: CombatAnim) -> void:
	combat_anim = anim
	_update_visual()
	# TODO: 动画播放完毕后重置
	await get_tree().create_timer(0.8).timeout
	if combat_anim == anim:
		combat_anim = CombatAnim.NONE
		_update_visual()


func _update_visual() -> void:
	# 选中高亮
	if highlight:
		match combat_state:
			CombatState.SELECTED:
				highlight.color = Color(1, 0.8, 0, 0.4)
				highlight.visible = true
			CombatState.ACTED:
				highlight.color = Color(0.3, 0.3, 0.3, 0.4)
				highlight.visible = true
			_:
				highlight.visible = false
	
	# 动画标签
	if state_label:
		match combat_anim:
			CombatAnim.ATTACKING:
				state_label.text = "攻击"
				state_label.visible = true
			CombatAnim.SKILL:
				state_label.text = "技能"
				state_label.visible = true
			CombatAnim.DAMAGED:
				state_label.text = "受伤"
				state_label.visible = true
			CombatAnim.DEATH:
				state_label.text = "死亡"
				state_label.visible = true
			_:
				state_label.visible = false
