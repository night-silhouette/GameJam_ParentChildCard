extends Node
# 【UI 展现层】只管渲染和玩家交互

@export var card_slot_scene: PackedScene 
@export var items_per_page: int = 4

@onready var grid_container: GridContainer = $"../card_vector"
@onready var btn_left = $"../切换/左切换"
@onready var btn_right= $"../切换/右切换"
var zone = Global.ZONE_CARD.BAG_ZONE

var current_page: int = 0


func _ready() -> void:
	# 核心：监听全局数据层的更新信号
	InventoryManager.bag_updated.connect(_on_global_data_updated)
	
	# 绑定翻页按钮
	btn_left.pressed.connect(_on_left_pressed)
	btn_right.pressed.connect(_on_right_pressed)
	SignalBus.left_clicked.connect(_on_card_left_clicked)
	SignalBus.right_clicked.connect(_on_card_right_clicked)
	
	# 初始化：如果打开 UI 时全局已经有数据了，直接刷一次
	_update_ui()


# 当局外数据层发生任何变化时，UI 自动被动响应
func _on_global_data_updated() -> void:
	_update_ui()


# 纯粹的渲染逻辑
func _update_ui() -> void:
	for child in grid_container.get_children():
		child.queue_free()
		
	# 💡 从全局单例中直接读取数据
	var global_list = InventoryManager.card_list
		
	if global_list.is_empty():
		btn_left.disabled = true
		btn_right.disabled = true
		return

	var start_idx = current_page * items_per_page
	var end_idx = min(start_idx + items_per_page, global_list.size())
	
	for i in range(start_idx, end_idx):
		var item_data = global_list[i]
		var slot_instance = card_slot_scene.instantiate()
		grid_container.add_child(slot_instance)
		
		if slot_instance.has_method("setup"):
			slot_instance.setup(item_data)
			
		if slot_instance.has_signal("left_clicked"):
			slot_instance.left_clicked.connect(_on_card_left_clicked)
		if slot_instance.has_signal("right_clicked"):
			slot_instance.right_clicked.connect(_on_card_right_clicked)
			
	btn_left.disabled = (current_page == 0)
	btn_right.disabled = (end_idx >= global_list.size())
	$"../页码".text = str(current_page+1)

func _on_left_pressed() -> void:
	if current_page > 0:
		current_page -= 1
		_update_ui()


func _on_right_pressed() -> void:
	var global_list = InventoryManager.card_list
	var max_page = (global_list.size() - 1) / items_per_page
	if current_page < max_page:
		current_page += 1
		_update_ui()


func _on_card_left_clicked(stuff_id: int) -> void:
	# 左键打开详情页（可以呼叫你的详细页 UI 或者是发射全局弹出信号）
	print("UI接收左键，准备展示详情，卡牌独一ID: ", stuff_id)


func _on_card_right_clicked(stuff_id: int,state) -> void:
	# 💡 右键点击时，UI 自己不删数据，而是命令全局数据层去改数据
	match state:
			0:
				InventoryManager.move_to_bag_zone(stuff_id)
			1:
				InventoryManager.move_to_sell_zone(stuff_id)
			2:
				InventoryManager.move_to_combat_zone(stuff_id)
