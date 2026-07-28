extends Node
@export var card_manager: Node
@export var zone : int = Global.ZONE_CARD.FREE_ZONE
@export var card : Control
@export var area : Area2D
var temp_id : int
var cards: Array   #
var is_drag : bool = false;
	

func _ready() -> void:
	# 监听数据层的变化，任何 Zone 的变动都会通过这里反映到 UI
	card_manager.UI_date_update.connect(refresh_ui)
	
func _process(delta: float) -> void:
	if not card: return
	
	if is_drag:
		# 1. 计算鼠标中心点目标位置
		var target_pos = card.get_global_mouse_position() - (card.size * card.scale / 2.0)
		# 2. 丝滑跟随 (20.0 是速度，值越大越粘手)
		card.global_position = card.global_position.lerp(target_pos, 20.0 * delta)
		
func refresh_ui():
	cards = card_manager.get_cards_by_zone(zone)
	_update_view()

func _update_view():
	if !cards.is_empty():
		card.import_card_data(cards[0].get("resouce"))
		temp_id = cards[0].get("temp_id")
		is_drag = true;
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


func _input(event: InputEvent) -> void:
	# 如果当前根本没有在拖拽，或者卡牌实例不存在，直接无视全局输入
	if not is_drag or not card: 
		return
		
	# 只要在拖拽状态下，无条件监听鼠标左键放开
	if event is InputEventMouseButton and event.button_index == MOUSE_BUTTON_LEFT:
		if not event.pressed: # 鼠标松开了
			
			# 在清空状态前，先把需要发送的数据存下来，防止下一步失效
			var temp_id_to_send = temp_id
			
			# 1. 立即重置拖拽状态，防止信号延迟导致二次触发
			is_drag = false 
			print(temp_id_to_send)
			# 2. 发送信号给后端/网络层
			SignalBus.exit_freecard.emit(temp_id_to_send)
			
			# 3. 拦截事件，防止这个松开事件穿透影响到地下的其他按钮
			get_viewport().set_input_as_handled()
