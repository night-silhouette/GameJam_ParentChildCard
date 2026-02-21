extends TextureRect
@export var range_y : int = 34;
@export var range_x : int = 32;
@export var speed : float = 100;
var mapWidth : int =1152;
var mapHeight : int = 640;

# Called when the node enters the scene tree for the first time.
var gap_y :int = 80;#这个就是最合适的距离和边框的距离。
var gap_x :int = 136;

var time_passed := 0.0;
var float_amplitude := 8.0;
var float_frequency := 2.0;

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
		

func _process(delta: float) -> void:
	pass;
	
func rand_normal(mean := 0.0, std_dev := 1.0):
	var u1 = randf()
	var u2 = randf()
	var z0 = sqrt(-2.0 * log(u1)) * cos(TAU * u2)
	return mean + z0 * std_dev
	
func change_rad():#改变角度
	match rotation_degrees :
		0:
			rotation_degrees = 90;
		90:
			rotation_degrees = 180;
		180:
			rotation_degrees = 270;
		270:
			rotation_degrees = 0;
		
func decide_change_rad():#判断改变方向
	if position.x <= 0 - size.x and rotation_degrees == 0 :
		change_rad();
	elif position.y <= 0 - size.x and rotation_degrees == 90 :
		change_rad();
	elif position.x >= mapWidth and rotation_degrees == 180 :
		change_rad();
	elif position.y >= mapHeight and rotation_degrees == 270 :
		change_rad();
		
		
	
