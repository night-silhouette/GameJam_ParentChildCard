extends Control

@export var loot_card_scene: PackedScene
@onready var gift_grid: GridContainer = $"容器/GridContainer"
@onready var loot_grid: GridContainer = $"牌库/card_vector"
@onready var submit_btn: TextureButton = $"万能按钮"


func _ready() -> void:
	InventoryManager.loot_updated.connect(refresh)
	refresh()
	submit_btn.button_down.connect(_on_submit)


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


func _on_submit() -> void:
	var card_list: Array = []
	for data in InventoryManager.loot_card_list:
		if data["zone"] == Global.ZONE_CARD.LOOT_ZONE:
			card_list.append({
				"card_id": data["card_id"],
				"stuff_id": data["stuff_id"]
			})
	SignalBus.request_loot_post.emit(card_list, InventoryManager.loot_id)
