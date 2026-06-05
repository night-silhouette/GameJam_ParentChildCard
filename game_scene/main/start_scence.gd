extends Node2D

@export var set_time = 2;

var sum_time = 0;

var token_success = false;

func _ready() -> void:
	SignalBus.token_validated_success.connect(_token_success)
	SignalBus.network_disconnected.connect(_token_fail)
	SignalBus.request_validate_token.emit();

# Called every frame. 'delta' is the elapsed time since the previous frame.
		
func _token_success() :
	token_success = true;
	# print("token成功")
func _token_fail():
	token_success = false
	# print("token错误")


func _on_timer_timeout() -> void:
	
	if token_success and Global.token_save :
		SignalBus.change_scence.emit("tomenu");
		SignalBus.change_ui.emit("tomenu");
		# print("menu")
	else :
		SignalBus.change_scence.emit("tologin");
		SignalBus.change_ui.emit("tologin")
		
