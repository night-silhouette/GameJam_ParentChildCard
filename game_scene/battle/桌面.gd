extends TextureRect

@export var card_manager: Node
@export var in_duration: float = 0.5

@export var cards_ui: Array[Control] = []
@export var zones: Array[int] = []
@export var hand: Sprite2D

var all_cards_data: Array = []

## 当前选中的己方牌信息
var _selected_temp_id: int = -1
var _selected_self_where: int = -1
## 选中的敌方目标位置
var _selected_opponent_where: int = -1:
	set(value):
		match value:
			0:
				hand_dir.flip_h = true;
			1:
				hand_dir.flip_h = false;
		_selected_opponent_where = value	
## 手部手势节点
@onready var hand_skill: Sprite2D = $"1"
@onready var hand_attack: Sprite2D = $"2"
@onready var hand_dir: TextureRect = $"指向"
@onready var btn_attack: TextureButton = $"战斗"
@onready var btn_skill: TextureButton = $"法术"
@onready var btn_reset: TextureButton = $"万能按钮"


func _ready() -> void:
	card_manager.UI_date_update.connect(refresh_ui)
	card_manager.combat_dto_changed.connect(_on_dto_changed)
	
	
	for card in cards_ui:
		if card:
			card.visible = false
	
	_connect_combat_signals()
	
	# 攻/法按钮：确认为选中+目标后的最终操作
	btn_attack.pressed.connect(_on_attack_pressed)
	btn_skill.pressed.connect(_on_skill_pressed)
	# 重置按钮：清空 DTO，返还能量
	btn_reset.pressed.connect(_on_reset_pressed)
	
	# 初始隐藏确认按钮
	_hide_confirm_buttons()
	SignalBus.enter_free.emit(hide_free)


func _connect_combat_signals() -> void:
	for card in cards_ui:
		if not card:
			continue
		if card.has_signal("combat_selected"):
			card.combat_selected.connect(_on_combat_card_selected)
		if card.has_signal("combat_target_changed"):
			card.combat_target_changed.connect(_on_combat_target_selected)


func refresh_ui():
	var loop_count = min(cards_ui.size(), zones.size())

	for i in range(loop_count):
		var current_card_ui = cards_ui[i]
		var current_zone = zones[i]

		if not current_card_ui:
			continue

		var zone_cards = card_manager.get_cards_by_zone(current_zone)

		if not zone_cards.is_empty():
			var icard = zone_cards[0]
			current_card_ui.update_card_data(icard)
			if not current_card_ui.visible or current_card_ui.modulate.a < 1.0:
				_fade_in(current_card_ui, in_duration)
		else:
			current_card_ui.visible = false
			


func _fade_in(target_card: Control, duration: float = 0.5) -> void:
	if not target_card: return
	
	target_card.visible = true
	target_card.modulate.a = 0.0
	target_card.process_mode = Node.PROCESS_MODE_INHERIT
	target_card.mouse_filter = Control.MOUSE_FILTER_STOP
	
	var tween = create_tween()
	tween.set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
	tween.tween_property(target_card, "modulate:a", 1.0, duration)


#region ==================== 战斗选中逻辑 ====================

## 己方牌被选中 → 显示指向手势，隐藏确认按钮
func _on_combat_card_selected(temp_id: int, self_where: int) -> void:
	if temp_id == -1:
		_reset_selection()
		return
	
	_selected_temp_id = temp_id
	_selected_self_where = self_where
	_selected_opponent_where = -1
	_hide_confirm_buttons()
	


## 点击敌方牌 → 记录目标，显示攻/法确认按钮
func _on_combat_target_selected(opponent_where: int) -> void:
	if _selected_temp_id == -1:
		return
	hand_dir.visible = true;
	_selected_opponent_where = opponent_where
	_show_confirm_buttons()


## 确认攻击
func _on_attack_pressed() -> void:
	if _selected_temp_id == -1 or _selected_opponent_where == -1:
		return
	_confirm_action(0)


## 确认法术
func _on_skill_pressed() -> void:
	if _selected_temp_id == -1 or _selected_opponent_where == -1:
		return
	_confirm_action(1)


## 执行确认：记录 DTO，标记已行动
func _confirm_action(behavior: int) -> void:
	var ok = card_manager.set_combat_dto(_selected_self_where, behavior, _selected_opponent_where, _selected_temp_id, {})
	if not ok:
		return  # 能量不足
	hand_dir.visible = false;
	# 标记己方牌为已行动
	var selected_card = _find_card_by_temp_id(_selected_temp_id)
	if selected_card and selected_card.has_method("mark_acted"):
		selected_card.mark_acted()
	
	_reset_selection()


## 重置按钮：清空所有 DTO，返还能量
func _on_reset_pressed() -> void:
	card_manager.clear_all_combat_dto()
	_reset_selection()


## DTO 变更回调：如果清空了，重置对应卡牌
func _on_dto_changed() -> void:
	var dp = card_manager.parent_combat_dto
	var dc = card_manager.child_combat_dto
	
	if dp.behavior == -1:
		_reset_card_by_self_where(0)
	if dc.behavior == -1:
		_reset_card_by_self_where(1)


func _reset_selection() -> void:
	_selected_temp_id = -1
	_selected_self_where = -1
	_selected_opponent_where = -1
	_hide_hand()
	_hide_confirm_buttons()


func _reset_card_by_self_where(self_where: int) -> void:
	for card in cards_ui:
		if card and card.get("self_where") == self_where and card.has_method("reset_combat"):
			card.reset_combat()


func _find_card_by_temp_id(temp_id: int):
	for card in cards_ui:
		if card and card.get("temp_id") == temp_id:
			return card
	return null


func _hide_hand() -> void:
	if hand_attack:
		hand_attack.visible = false
	if hand_skill:
		hand_skill.visible = false
	if hand_dir:
		hand_dir.visible = false


func _show_confirm_buttons() -> void:
	if btn_attack:
		btn_attack.visible = true
	if btn_skill:
		btn_skill.visible = true


func _hide_confirm_buttons() -> void:
	if btn_attack:
		btn_attack.visible = false
	if btn_skill:
		btn_skill.visible = false


func reset_all_combat_cards() -> void:
	for card in cards_ui:
		if card and card.has_method("reset_combat"):
			card.reset_combat()
	_reset_selection()
#endregion

func hide_free() :
	_hide_confirm_buttons()
