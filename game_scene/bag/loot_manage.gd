extends Control
# 挂载在"获取战利品"节点上，管理排库(LOOT_ZONE)和礼物区(GIFT_ZONE)

@export var loot_card_scene: PackedScene
@onready var gift_grid: GridContainer = $gift_grid
@onready var loot_grid: GridContainer = $loot_grid


func _ready() -> void:
	InventoryManager.loot_updated.connect(refresh)
	refresh()


func refresh() -> void:
	_populate(gift_grid, Global.ZONE_CARD.GIFT_ZONE)
	_populate(loot_grid, Global.ZONE_CARD.LOOT_ZONE)


func _populate(grid: GridContainer, zone: int) -> void:
	for child in grid.get_children():
		child.queue_free()

	for data in InventoryManager.loot_card_list:
		if data["zone"] != zone:
			continue
		var card = loot_card_scene.instantiate()
		grid.add_child(card)
		if card.has_method("setup"):
			card.setup(data)
