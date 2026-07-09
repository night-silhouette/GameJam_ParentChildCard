extends TextureRect
@export var card_manager: Node
@export var card : Control
@export var in_duration :	float
var cards: Array  
# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	card_manager.UI_date_update.connect(refresh_ui)
	card.visible = false;
	
func refresh_ui():
	cards = card_manager.get_cards_by_zone(Global.ZONE_CARD.SPELL_ZONE)
	if !cards.is_empty():
		var icard = cards[0];
		card.update_card_data(icard);
		if icard.get("need_operate", false):
			card.enter_need_operate()
		else:
			card.exit_need_operate()
		_fade_in(in_duration)
	else:
		card.visible = false;
## @duration: 浮现持续的时间（秒）
func _fade_in(duration: float = 0.5) -> void:
	if not card: return
	
	# 1. 确保基础状态正确
	card.visible = true
	card.modulate.a = 0.0  # 先变完全透明
	card.process_mode = Node.PROCESS_MODE_INHERIT
	card.mouse_filter = Control.MOUSE_FILTER_STOP
	
	# 2. 创建并执行动画
	var tween = create_tween()
	tween.set_trans(Tween.TRANS_QUAD).set_ease(Tween.EASE_OUT)
	tween.tween_property(card, "modulate:a", 1.0, duration)


		
