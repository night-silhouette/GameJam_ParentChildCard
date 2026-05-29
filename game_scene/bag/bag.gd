extends Control


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	SignalBus.request_bag_card.emit();
	SignalBus.request_get_self_gold.emit()
	SignalBus.left_clicked.connect(_left_clicked)
	SignalBus.right_clicked.connect(_right_clicked)
func _left_clicked():
	pass;
	
func _right_clicked(stuff_id):
	InventoryManager.move_to_sell_zone(stuff_id)
	SignalBus.request_bag_card.emit();
	SignalBus.request_get_self_gold.emit()
