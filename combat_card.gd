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
#endregion


func _ready() -> void:
	super._ready()
	
	# 根据 zone 判断是否为己方牌
	if zone in [Global.ZONE_CARD.PARENT_BATTLE_ZONE, Global.ZONE_CARD.CHILD_BATTLE_ZONE]:
		is_own = true
		self_where = 0 if zone == Global.ZONE_CARD.PARENT_BATTLE_ZONE else 1
	
	# 监听 HP 变化动画
	SignalBus.ani_hp_change_enter.connect(_on_hp_change_enter)


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


#region HP 变化数值动画
const HP_COLOR = {
	0: Color(0.2, 0.9, 0.2),   # Damage → 绿
	1: Color(1.0, 0.2, 0.2),   # Heal → 红
	2: Color(1.0, 1.0, 1.0),   # TrueDamage → 白
}

## 收到 ani_hp_change_enter 信号，匹配 temp_id 决定是否播放
func _on_hp_change_enter(temp_id: int, category: int, value: int) -> void:
	if temp_id != self.temp_id:
		return
	_play_value_change(category, value)


## 播放数值变化动画：浮现 → 放大回缩 → 淡出消失
func _play_value_change(category: int, value: int) -> void:
	var color = HP_COLOR.get(category, Color.WHITE)
	
	# 创建数值 Label
	var label = Label.new()
	label.text = ("+" if category == 1 else "-") + str(abs(value))
	label.label_settings = LabelSettings.new()
	label.label_settings.font_size = 40
	label.label_settings.font_color = color
	label.horizontal_alignment = HORIZONTAL_ALIGNMENT_CENTER
	label.vertical_alignment = VERTICAL_ALIGNMENT_CENTER
	label.position = Vector2(size.x / 2 - 50, -(size.y * 0.3))
	label.size = Vector2(100, 50)
	label.scale = Vector2.ZERO
	label.modulate.a = 1.0
	add_child(label)
	
	var tween = create_tween()
	tween.set_parallel(true)
	
	# 浮现上移
	tween.tween_property(label, "position:y", label.position.y - 40, 0.5)
	# 放大
	tween.tween_property(label, "scale", Vector2(1.5, 1.5), 0.15)
	# 回缩
	tween.tween_property(label, "scale", Vector2.ONE, 0.15).set_delay(0.15)
	# 淡出
	tween.tween_property(label, "modulate:a", 0.0, 0.4).set_delay(0.2)
	
	# 动画结束后清理
	tween.chain().tween_callback(label.queue_free)
#endregion
