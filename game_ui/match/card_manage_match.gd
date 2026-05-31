extends Node
# 【UI 展现层】只管渲染和玩家交互

@export var items_per_page: int = 5

@onready var bag_card: GridContainer = $"../牌库底板/bag"
@onready var match_card: GridContainer = $"../GridContainer"
@onready var btn_left = $"../牌库底板/切换/左切换"
@onready var btn_right= $"../牌库底板/切换/右切换"
@onready var page_label = $"../页码"

var current_page: int = 0

# 缓存场景中已经添加好的卡槽节点
var bag_slots: Array = []
var match_slots: Array = []

func _ready() -> void:
	# 1. 获取已经在编辑器里摆好的 5 个节点
	bag_slots = bag_card.get_children()
	match_slots = match_card.get_children()
	
	# 2. 核心：监听全局数据层的更新信号
	InventoryManager.bag_updated.connect(_update_ui)
	InventoryManager.bag_updated.connect(_update_ui)
	
	# 3. 绑定翻页按钮与全局点击信号
	btn_left.pressed.connect(_on_left_pressed)
	btn_right.pressed.connect(_on_right_pressed)
	SignalBus.left_clicked.connect(_on_card_left_clicked)
	SignalBus.right_clicked.connect(_on_card_right_clicked)
	
	# 4. 初始化：刷一次数据
	_update_ui()

# 当局外数据层发生任何变化时，UI 自动被动响应
func _update_ui() -> void:
	_update_bag_zone()
	_update_match_zone()

# 专门处理背包区的刷新（带翻页逻辑）
func _update_bag_zone() -> void:
	# 💡 从全局单例获取背包数据列表
	var bag_list = InventoryManager.get_cards_in_zone(Global.ZONE_CARD.BAG_ZONE) # 需在 InventoryManager 中实现此方法返回背包数组
	var total_items = bag_list.size()
	
	# 更新翻页按钮状态和页码文本
	var max_page = max(0, (total_items - 1) / items_per_page)
	# 如果当前页超出了最大页（比如卖掉东西后总数减少），自动退回有效页
	if current_page > max_page:
		current_page = max_page
		
	btn_left.disabled = (current_page == 0)
	btn_right.disabled = (current_page >= max_page) or total_items == 0
	page_label.text = str(current_page + 1)
	
	var start_idx = current_page * items_per_page
	
	# 将数据映射到预设好的 5 个卡槽上
	for i in range(bag_slots.size()):
		var slot = bag_slots[i]
		var data_idx = start_idx + i
		
		if data_idx < total_items:
			slot.visible = true
			if slot.has_method("setup"):
				slot.setup(bag_list[data_idx])
		else:
			slot.visible = false

# 专门处理出战区的刷新（无翻页，直接映射前5个）
func _update_match_zone() -> void:
	# 💡 从全局单例获取出战区数据列表
	var match_list = InventoryManager.get_cards_in_zone(Global.ZONE_CARD.MATCH_ZONE) # 需在 InventoryManager 中实现此方法返回出战数组
	var total_items = match_list.size()
	
	# 将数据映射到预设好的 5 个卡槽上
	for i in range(match_slots.size()):
		var slot = match_slots[i]
		
		if i < total_items:
			slot.visible = true
			if slot.has_method("setup"):
				slot.setup(match_list[i])
		else:
			slot.visible = false

# --- 交互操作 ---

func _on_left_pressed() -> void:
	if current_page > 0:
		current_page -= 1
		_update_bag_zone() # 只刷新背包区即可

func _on_right_pressed() -> void:
	var bag_list = InventoryManager.get_bag_list()
	var max_page = max(0, (bag_list.size() - 1) / items_per_page)
	
	if current_page < max_page:
		current_page += 1
		_update_bag_zone() # 只刷新背包区即可

func _on_card_left_clicked(stuff_id: int) -> void:
	InventoryManager.move_to_combat_zone(stuff_id)

func _on_card_right_clicked(stuff_id: int, state: int) -> void:
	# 💡 右键点击时，UI 自己不改数据，而是命令全局数据层去改
	match state:
		0:
			InventoryManager.move_to_bag_zone(stuff_id)
		1:
			InventoryManager.move_to_sell_zone(stuff_id)
		2:
			InventoryManager.move_to_combat_zone(stuff_id)
