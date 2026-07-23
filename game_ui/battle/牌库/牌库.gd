extends Control

@export var object_pool : Node
@export var card_vector: GridContainer
@export var card_manager: Node


var cards: Array   # 你的全部牌
var page_size: int = 5
var start_index: int = 0

func _ready() -> void:
	# 监听数据层的变化，任何 Zone 的变动都会通过这里反映到 UI
	card_manager.change_card_zone.connect(refresh_ui)
	card_manager.UI_date_update.connect(refresh_ui)
	# 1. 初始化预填充
	for i in range(page_size):
		var child = object_pool.get_card()
		child.hide()
		card_vector.add_child(child)
		
		## 连接子节点的信号：当子节点被玩家“抓起”时
		## 注意：子节点只需传出它持有的 card_data 即可
		if child.has_signal("request_drag"):
			child.request_drag.connect(_on_card_request_drag)

func refresh_ui():
	# 从数据层重新拉取当前 Zone 的所有卡牌数据
	cards = card_manager.get_cards_by_zone(Global.ZONE_CARD.DECK_ZONE)
	_update_view()

func _update_view():
	var page_cards = get_current_page()
	var ui_nodes = card_vector.get_children()

	for i in range(ui_nodes.size()):
		var child = ui_nodes[i]
		if i < page_cards.size():
			var card_data = page_cards[i]
			child.update_card_data(card_data)
			child.show()
			# 如果该卡有 FREE_ZONE 副本，显示半透明占位
			if card_manager.has_free_zone_duplicate(card_data.get("temp_id")):
				child.modulate.a = 0.3
			else:
				child.modulate.a = 1.0
			# 根据 card_data 的 need_operate 标记切换状态
			# 跳过正在 HOVERED 的卡牌，避免覆盖 hover 状态导致 mouse_exited 失效
			if child.current_state != child.CardState.HOVERED:
				if card_data.get("need_operate", false):
					child.enter_need_operate(true)
				else:
					child.exit_need_operate(true)
		else:
			child.hide()

func get_current_page() -> Array:
	var end_index = min(start_index + page_size, cards.size())
	return cards.slice(start_index, end_index)

# ================= 核心：配合“游荡对象”的逻辑 =================

func _on_card_request_drag(card_data):
	# 原卡保留在 DECK_ZONE 作为半透明占位，FREE_ZONE 副本由 enter_freecard 信号创建
	pass

# ================= 翻页逻辑（保持不变） =================

func next_page():
	if start_index + page_size < cards.size():
		start_index += page_size
		refresh_ui()

func prev_page():
	start_index = max(start_index - page_size, 0)
	refresh_ui()

func _on_右切换_button_down() -> void:
	next_page()

func _on_左切换_button_down() -> void:
	prev_page()
