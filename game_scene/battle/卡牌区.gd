extends Control

@export var object_pool : Node
@export var card_vector: GridContainer
@export var card_manager: Node
@export var page_size: int = 1
@export var zone : int
var cards: Array   # 你的全部牌

var start_index: int = 0

func _ready() -> void:
	# 监听数据层的变化，任何 Zone 的变动都会通过这里反映到 UI
	card_manager.UI_date_update.connect(refresh_ui)
	
	# 1. 初始化预填充
	for i in range(page_size):
		var child = object_pool.get_card()
		child.hide()
		card_vector.add_child(child)
		
		
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
			child.update_card(page_cards[i])
			child.show()
			# 确保卡牌视觉状态正常（比如之前被拖拽时可能变透明了，这里重置一下）
			child.modulate.a = 1.0 
		else:
			child.hide()

func get_current_page() -> Array:
	var end_index = min(start_index + page_size, cards.size())
	return cards.slice(start_index, end_index)

# ================= 核心：配合“游荡对象”的逻辑 =================

func _on_card_request_drag(card_data):
	# 1. 激活全局游荡对象（DragProxy）
	# 假设你有一个全局单例 DragManager 
	#DragManager.start_drag(card_data)
	
	# 2. 核心操作：修改数据层！
	# 告诉 card_manager：这张牌现在离开 DECK_ZONE 了，进入临时状态
	# 只要这一步执行并触发了 show_card 信号，refresh_ui 就会自动执行
	# 从而实现“瞬间补位”
	card_manager.remove_card_from_view(card_data)
