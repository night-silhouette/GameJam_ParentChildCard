extends Node2D

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	ScenceManage.register_root($Scence_Root);
	SignalBus.change_scence.emit("start");
	UiManage.register_root($UI_Root);
	
