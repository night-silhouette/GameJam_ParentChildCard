extends Node2D
@export var token_save :bool = false;
# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	ScenceManage.register_root($Scence_Root);#注册场景节点
	
	UiManage.register_root($UI_Root);#注册UI节点
	
	SignalBus.change_scence.emit("start");
	
	Global.token_save = token_save;
	
	
