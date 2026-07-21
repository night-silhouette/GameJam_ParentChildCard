extends Control


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	
	SignalBus.match_success.connect(_match_success);
	SignalBus.soft_reconnect.connect(_soft_reconnect)
	$"标识".txt = "Click"
	$"游戏开始".txt = "匹配！"
	$"牌库".txt = "仓库！！"
	
	
func _match_success(t):
	
	Global.init_battle_time = t;
	SignalBus.change_scence.emit("tobattle");
	SignalBus.change_ui.emit("tomenu");
	
func _soft_reconnect():
	SignalBus.change_scence.emit("tobattle");
	SignalBus.change_ui.emit("tomenu");
	SignalBus.request_reconnect_query.emit()
