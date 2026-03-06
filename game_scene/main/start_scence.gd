extends Node2D

@export var set_time = 3;

var sum_time = 0;

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
		sum_time += delta;
		if sum_time > set_time : 
			SignalBus.change_scence.emit("tologin");
			SignalBus.change_ui.emit("tologin")
		
		
		
