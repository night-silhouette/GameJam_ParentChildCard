extends TextureRect

@export var card_manager: Node
@export var in_duration: float = 0.5

# 1. 使用数组将 4 个 Card 节点和 4 个 Zone 顺序列出来
# 请在编辑器中把对应的节点和 Zone 按相同的顺序拖入/填入数组中
@export var cards_ui: Array[Control] = []
@export var zones: Array[int] = []

# 用来存储从数据层获取的原始卡牌数据（如果需要的话）
var all_cards_data: Array = [] 

func _ready() -> void:
	card_manager.UI_date_update.connect(refresh_ui)
	
	# 初始时将所有卡片隐藏
	for card in cards_ui:
		if card:
			card.visible = false

func refresh_ui():
	# 确保两个数组长度一致，防止越界崩溃
	var loop_count = min(cards_ui.size(), zones.size())
	
	for i in range(loop_count):
		var current_card_ui = cards_ui[i]
		var current_zone = zones[i]
		
		# 安全检查，防止编辑器里漏拖节点
		if not current_card_ui: 
			continue
			
		# 获取当前 Zone 的卡牌数组
		var zone_cards = card_manager.get_cards_by_zone(current_zone)
		
		# 如果该 Zone 有牌，取第一张（index 0）进行更新并淡入
		if not zone_cards.is_empty():
			var icard = zone_cards[0]
			current_card_ui.update_card_data(icard)
			
			# 只有当卡片之前是隐藏状态，或者已经完全透明时，才触发淡入动画
			# 这样可以避免数据频繁更新时，动画不停地从头播放导致闪烁
			if not current_card_ui.visible or current_card_ui.modulate.a < 1.0:
				_fade_in(current_card_ui, in_duration)
		else:
			# 如果该 Zone 没牌了，直接隐藏
			current_card_ui.visible = false

## 传入指定的 card 节点进行淡入
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
