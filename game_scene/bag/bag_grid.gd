extends Control

# 单个卡牌槽的 UI 预制体（需要在检查器中拖入你的 CardSlot.tscn）
@export var card_slot_scene: PackedScene 
@export var items_per_page: int = 20 # 每页显示的卡牌数量
@export_dir var cards_folder_path: String = "res://resources/cards"

@onready var grid_container: GridContainer = $GridContainer
@onready var btn_left: Button = $LeftButton
@onready var btn_right: Button = $RightButton

# 当前所在页码
var current_page: int = 0

# 核心需求：利用 Setter 监听数据变化
var card_list: Array = []:
	set(value):
		card_list = value
		current_page = 0 # 拿到新背包数据时，重置回第一页
		_update_ui() # 自动触发 UI 刷新


func _ready() -> void:
	# 绑定网络/系统信号
	SignalBus.get_card_bag.connect(_get_card_bag)
	
	# 绑定按钮点击事件
	btn_left.pressed.connect(_on_left_pressed)
	btn_right.pressed.connect(_on_right_pressed)


# 接收原始数据并处理成包含 Resource 的完整字典列表
func _get_card_bag(incoming_card_data: Array):
	var temp_list = []
	
	for item in incoming_card_data:
		var card_res = _find_card_resource_by_id(item["card_id"])
		
		var data = {
			"stuff_id": item["stuff_id"],
			"card_id": item["card_id"],
			"price": item["price"],
			"zone": Global.ZONE_CARD.BAG_ZONE,
			"resource": card_res
		}
		temp_list.append(data)
		
	# 注意：这里必须整体重新赋值，才能触发 set(value) 监听器
	card_list = temp_list


# 核心渲染函数：根据当前页码刷新 Grid 内部的卡牌
func _update_ui() -> void:
	# 1. 清空当前网格内所有的旧 UI 节点
	for child in grid_container.get_children():
		child.queue_free()
		
	if card_list.is_empty():
		btn_left.disabled = true
		btn_right.disabled = true
		return

	# 2. 计算当前页需要截取的数组索引范围
	var start_idx = current_page * items_per_page
	var end_idx = min(start_idx + items_per_page, card_list.size())
	
	# 3. 实例化并填充卡牌组件
	for i in range(start_idx, end_idx):
		var item_data = card_list[i]
		var slot_instance = card_slot_scene.instantiate()
		grid_container.add_child(slot_instance)
		
		# 假设你的单个卡牌 UI 脚本里写了一个用于接收数据的初始化函数，例如 func setup(data: Dictionary)
		if slot_instance.has_method("setup"):
			slot_instance.setup(item_data)
			
	# 4. 动态更新翻页按钮的禁用状态
	btn_left.disabled = (current_page == 0)
	btn_right.disabled = (end_idx >= card_list.size())


# 左翻页
func _on_left_pressed() -> void:
	if current_page > 0:
		current_page -= 1
		_update_ui()


# 右翻页
func _on_right_pressed() -> void:
	var max_page = (card_list.size() - 1) / items_per_page
	if current_page < max_page:
		current_page += 1
		_update_ui()


# 利用路径拼接直接载入 card_X.tres
func _find_card_resource_by_id(target_id: int) -> CardResource:
	var resource_path = cards_folder_path + "/card_" + str(target_id) + ".tres"
	if ResourceLoader.exists(resource_path):
		return load(resource_path) as CardResource
	push_warning("未找到卡牌资源文件: " + resource_path)
	return null


# 根据 stuff_id 修改 zone 为 sell_zone
func move_to_sell_zone(stuff_id: int) -> bool:
	for item in card_list:
		if item["stuff_id"] == stuff_id:
			item["zone"] = Global.ZONE_CARD.SELL_ZONE
			
			# 特别注意：在 Godot 中，仅修改数组内部字典的键值，并不会触发 Array 变量本身的 Setter。
			# 所以这里我们需要手动调用一次 UI 刷新，以同步最新的 Zone 状态。
			_update_ui() 
			return true
	return false
