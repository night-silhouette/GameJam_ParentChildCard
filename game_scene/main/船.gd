extends TextureRect
@export var range_y : int = 34;
@export var range_x : int = 32;
var mapWidth : int =1152;
var mapHeight : int = 640;

# Called when the node enters the scene tree for the first time.
var gap_y :int = 80;#这个就是最合适的距离和边框的距离。
var gap_x :int = 136;


func _ready() -> void:
	var random_x = [gap_x,mapWidth-gap_x];
	var random_y = [gap_y,mapHeight-gap_y];
	
	if randf() < 0.5 : #是在两侧；结果 = 值1 if 条件 else 值2
		
		position.x = random_x.pick_random();
		rotation = deg_to_rad(90) if position.x == gap_x else deg_to_rad(270);
		position.y = randf_range(0-size.x,mapHeight);
		
	else : #在顶部和底部
		
		position.y = random_y.pick_random();
		rotation = deg_to_rad(0) if position.y == gap_y else deg_to_rad(180);
		position.x = randf_range(0-size.x,mapWidth);

# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass
