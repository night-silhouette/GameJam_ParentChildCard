extends Control


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	SignalBus.match_success.connect(_match_success);


func _match_success(t):
	
	Global.init_battle_time = t;
	
	SignalBus.change_scence.emit("tobattle");	
	SignalBus.change_ui.emit("tobattle");	
	
