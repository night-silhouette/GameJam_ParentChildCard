extends Control
func _ready() -> void:
	SignalBus.request_get_self_cards.emit();
	
func _self_cards_updated(data):
	print(data);
	
