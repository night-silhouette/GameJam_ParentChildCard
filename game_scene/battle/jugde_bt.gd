extends Node

@export var texture_state_0 : Texture2D
@export var texture_state_1 : Texture2D
@export var texture_state_2 : Texture2D

@export var anim_fade_duration: float = 0.3
@export var anim_scale_duration: float = 0.35

@onready var texture_rect_1 : TextureRect = $TextureRect
@onready var texture_rect_2 : TextureRect = $TextureRect2

var tween_1 : Tween
var tween_2 : Tween

var _pending_anim_count: int = 0

func _on_anim_finished():
	_pending_anim_count -= 1
	if _pending_anim_count <= 0:
		_pending_anim_count = 0
		SignalBus.request_end_animation.emit()

var judge_data : Array = []:
	set(new_value):
		judge_data = new_value
		_refresh_ui_by_data()

func _ready() -> void:
	_refresh_ui_by_data()


func _refresh_ui_by_data() -> void:
	if not is_inside_tree():
		await ready

	if judge_data.size() < 2:
		_hide_node_instantly(texture_rect_1)
		_hide_node_instantly(texture_rect_2)
		return

	var val_1 : int = judge_data[0]
	var val_2 : int = judge_data[1]

	_pending_anim_count = 0

	_update_node_state(texture_rect_1, val_1, 1)
	_update_node_state(texture_rect_2, val_2, 2)

	if _pending_anim_count == 0:
		SignalBus.request_end_animation.emit()


func update_single_judge_data(index: int, value: int) -> void:
	while judge_data.size() < 2:
		judge_data.append(-1)

	judge_data[index] = value

	_refresh_ui_by_data()


func _update_node_state(target_rect: TextureRect, value: int, index: int) -> void:
	if not is_instance_valid(target_rect):
		return

	if value == -1:
		_hide_node_instantly(target_rect)
		return

	if index == 1 and tween_1: tween_1.kill()
	if index == 2 and tween_2: tween_2.kill()

	match value:
		0: target_rect.texture = texture_state_0
		1: target_rect.texture = texture_state_1
		2: target_rect.texture = texture_state_2
		_:
			_hide_node_instantly(target_rect)
			return
	if index == 2:
		target_rect.rotation = PI

	if not target_rect.visible or target_rect.modulate.a < 0.1:
		target_rect.pivot_offset = target_rect.size / 2

		target_rect.modulate.a = 0.0
		target_rect.scale = Vector2(0.8, 0.8)
		target_rect.visible = true

		var new_tween = create_tween().set_parallel(true)

		new_tween.tween_property(target_rect, "modulate:a", 1.0, anim_fade_duration)\
			.set_trans(Tween.TRANS_SINE).set_ease(Tween.EASE_OUT)

		new_tween.tween_property(target_rect, "scale", Vector2.ONE, anim_scale_duration)\
			.set_trans(Tween.TRANS_BACK).set_ease(Tween.EASE_OUT)

		if index == 1: tween_1 = new_tween
		if index == 2: tween_2 = new_tween

		_pending_anim_count += 1
		new_tween.finished.connect(_on_anim_finished)
	else:
		target_rect.visible = true
		target_rect.modulate.a = 1.0
		target_rect.scale = Vector2.ONE


func _hide_node_instantly(target_rect: TextureRect) -> void:
	if is_instance_valid(target_rect):
		target_rect.visible = false
		target_rect.modulate.a = 0.0
		target_rect.scale = Vector2(0.8, 0.8)
