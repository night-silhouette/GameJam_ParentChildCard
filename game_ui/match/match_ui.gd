extends Control


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	SignalBus.request_bag_card.emit()


# Called every frame. 'delta' is the elapsed time since the previous frame.


func _on_街机按钮_button_down() -> void:
	SignalBus.to_connect_ws.emit()


func _on_返回_button_down() -> void:
	SignalBus.change_ui.emit("tomenu")
