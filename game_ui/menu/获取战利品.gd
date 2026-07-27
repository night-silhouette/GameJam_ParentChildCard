extends Control

@export var loot_card_scene: PackedScene
@onready var gift_grid: GridContainer = $"容器/GridContainer"
@onready var loot_grid: GridContainer = $"牌库/GridContainer"
@onready var submit_btn: TextureButton = $"万能按钮"


func _ready() -> void:
	InventoryManager.loot_updated.connect(refresh)
	refresh()
	submit_btn.button_down.connect(_on_submit)


func refresh() -> void:
	_populate(gift_grid, Global.ZONE_CARD.GIFT_ZONE)
	_populate(loot_grid, Global.ZONE_CARD.LOOT_ZONE)


const CARD_COLUMNS := 5

func _make_placeholder() -> Control:
	var p = loot_card_scene.instantiate()
	p.modulate.a = 0.0
	p.mouse_filter = Control.MOUSE_FILTER_IGNORE
	p.process_mode = Node.PROCESS_MODE_DISABLED
	return p


func _populate(grid: GridContainer, zone: int) -> void:
	for child in grid.get_children():
		child.queue_free()

	var real_cards: Array = []
	for data in InventoryManager.loot_card_list:
		if data["zone"] != zone:
			continue
		real_cards.append(data)

	# 第一排：5张透明占位（顶部hover空间）
	for _i in range(CARD_COLUMNS):
		grid.add_child(_make_placeholder())

	# 实际卡牌
	for data in real_cards:
		var card = loot_card_scene.instantiate()
		grid.add_child(card)
		if card.has_method("setup"):
			card.setup(data)

	# 底部填充：补满当前行 + 额外一整排（底部hover空间）
	var total := CARD_COLUMNS + real_cards.size()
	var fill := (CARD_COLUMNS - total % CARD_COLUMNS) % CARD_COLUMNS
	fill += CARD_COLUMNS
	for _i in range(fill):
		grid.add_child(_make_placeholder())


func _on_submit() -> void:
	var card_list: Array = []
	for data in InventoryManager.loot_card_list:
		if data["zone"] == Global.ZONE_CARD.LOOT_ZONE:
			card_list.append(data["card_id"])
	print(card_list)
	SignalBus.request_loot_post.emit(card_list, InventoryManager.loot_id)
