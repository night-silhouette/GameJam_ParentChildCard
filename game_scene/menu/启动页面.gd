extends Node2D
signal change_scence(address : String)
@export var set_time = 10;
@export var next_scence_address =  "res://game_scene/main/登录页面.tscn"
var sum_time = 0;

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	pass # Replace with function body.


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	sum_time += delta;
	if sum_time > set_time :
		change_scence.emit(next_scence_address);
		
		
		
