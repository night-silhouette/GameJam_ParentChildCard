extends Node
@export var card_manager: Node
@export var page_size: int = 1
@export var zone : int
@export var card : Control
@export var area : Area2D
var cards: Array   # 你的全部牌


func _ready() -> void:
	# 监听数据层的变化，任何 Zone 的变动都会通过这里反映到 UI
	card_manager.UI_date_update.connect(refresh_ui)
	
		
func refresh_ui():
	cards = card_manager.get_cards_by_zone(zone)
	_update_view()

func _update_view():
	if !	cards.is_empty():
		card.update_card_data(cards[0])
		_activate();
	else:
		_deactivate();

func _deactivate():

	# 隐藏
	card.visible = false

	# 停止 process
	card.process_mode = Node.PROCESS_MODE_DISABLED

	# 禁止输入
	card.mouse_filter = Control.MOUSE_FILTER_IGNORE


	area.monitoring = false
	area.monitorable = false

func _activate():
	# 显示
	card.visible = true

	# 恢复 process
	card.process_mode = Node.PROCESS_MODE_INHERIT

	# 恢复输入
	card.mouse_filter = Control.MOUSE_FILTER_STOP
	area.monitoring = true
	area.monitorable = true
