extends TextureRect
@export var range_y : int = 34;
@export var range_x : int = 32;
@export var speed : float = 100;
@export var amplitude : float = 30;
@export var ro_speed := 0.0;
@export var noise_frequent = 0.1;

var noise := FastNoiseLite.new();
@export var time := 0.0
var base_pos : Vector2

var mapWidth : int =1152;
var mapHeight : int = 640;

#define一下
const left =1;
const right = 2;
const down = 3;
const up = 4;

# Called when the node enters the scene tree for the first time.
var gap_y :int = 140; #这个就是最合适的距离和边框的距离，注。
var gap_x :int = 140;

@export var locate_status : int = 0;

var x_min = gap_x-size.x;
var x_max = mapWidth-gap_x;
var y_min = gap_y-size.x;
var y_max = mapHeight-gap_y;



func _ready() -> void:		
	random_position();
	process_noise();
	
func _process(delta: float) -> void:
#速度需要发生变化。上下的浮动同时也要发生变化。那自然而然的角度应该也有变化，这就说明了直接用角度去判断撞墙掉头是非常麻烦的。会出很多bug。
#已经确定了隐形的直线，和一开始说的那样，那就用一个数组去处理状态，直接用int吗？1234四个状态，你用数组本身也要去读取吧？可是数组的上限更高吧，那还不如直接就用一个整数吧，我说实在不行可以当作指针。
	decide_change_rad();
	process_move(delta);
	
func change_rad_position():#转角处的变化
	match locate_status :
		down:
			locate_status = left;
			rotation_degrees = 90;
			position = Vector2(0+(gap_x - size.x),mapHeight+size.x/2)
		left:
			locate_status = up;
			rotation_degrees = 180;
			position = Vector2(0-size.x,y_min);
		up:
			locate_status = right ; 
			rotation_degrees = 270;
			position = Vector2(mapWidth-gap_x,0-size.x/2);
		right:
			locate_status = down;
			rotation_degrees = 0;
			position = Vector2(mapWidth+size.x/2,mapHeight - gap_y);

func decide_change_rad():#判断改变方向
	if position.x <= 0 - size.x and locate_status == down :
		change_rad_position();
	elif position.y <= 0 - size.x and locate_status == left :
		change_rad_position();
	elif position.x >= mapWidth  and locate_status == up :
		change_rad_position();
	elif position.y >= mapHeight  and locate_status == right :
		change_rad_position();

func process_noise(): #噪音的处理
	noise.seed = randi();
	noise.frequency = noise_frequent ;
	base_pos = position;

func process_move(delta):#运动处理

	time += delta;
	var n = noise.get_noise_1d(time);
	var target_position :Vector2;
	match locate_status :
		down:
			target_position.x = position.x - speed;
			target_position.y = y_max + n * amplitude ;
			
		left:
			target_position.y = position.y - speed;
			target_position.x = x_min + n * amplitude ;
		up:
			target_position.x = position.x + speed;
			target_position.y = y_min + n * amplitude;
		right:
			target_position.y = speed + position.y;
			target_position.x = x_max + n * amplitude;
	var angle = position.angle_to_point(target_position);
	rotation = lerp(rotation,angle,ro_speed * delta);
	position = target_position;

	
func random_position():#随机出现

	
	var random_x = [x_min,x_max];
	var random_y = [y_min,y_max];
	
	if randf() < 0.5 : #是在两侧；结果 = 值1 if 条件 else 值2
	
		position.x = random_x.pick_random();
		
		if position.x == x_min :
			rotation = deg_to_rad(90); 
			locate_status = left;
		else :
			rotation = deg_to_rad(270);
			locate_status = right;
		position.y = randf_range(0,mapHeight-size.x);
	
	else : #在顶部和底部
		position.y = random_y.pick_random();
		
		if position.y == y_max :
			rotation = deg_to_rad(0)
			locate_status = down;
		else:
			rotation = deg_to_rad(180);
			locate_status = up;
			
		position.x = randf_range(0,mapWidth-size.x);
		
