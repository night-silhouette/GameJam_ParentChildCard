extends Control

@export var loot_card_scene: PackedScene
@onready var gift_grid: GridContainer = $"容器/GridContainer"
@onready var loot_grid: GridContainer = $"牌库/card_vector"


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
