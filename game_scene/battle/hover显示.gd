extends Control

@export var card_manager: Node

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	card_manager.UI_date_update.connect(_refresh_ui)


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _refresh_ui():
	var _temp_id = card_manager.hover_card;
	if _temp_id == -1:
		visible = false;
		
	else :
		var card = card_manager.select_card_by_key(_temp_id,"temp_id")
		$hp.text = str(card.get("hp", ""))
		$damage.text = str(card.get("damage", ""))
		$"详细".text = str(card.get("spell_des", ""))
		$Control.update_card_data(card)
		_fade_in(self);
		
func _fade_in(target_card: Control, duration: float = 0.5) -> void:
	if not target_card: return
	
	# 1. 确保基础状态正确
	target_card.visible = true
	target_card.modulate.a = 0.0  # 先变完全透明
	target_card.process_mode = Node.PROCESS_MODE_INHERIT
	target_card.mouse_filter = Control.MOUSE_FILTER_STOP
	
	# 2. 创建并执行动画
	var tween = create_tween()
	tween.set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
	tween.tween_property(target_card, "modulate:a", 1.0, duration)
