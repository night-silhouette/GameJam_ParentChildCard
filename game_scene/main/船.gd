extends TextureRect
@export var range_y : int = 34;
@export var range_x : int = 32;
@export var speed : float = 100;
@export var amplitude : float = 30;

var noise := FastNoiseLite.new();
var time := 0

var mapWidth : int =1152;
var mapHeight : int = 640;

#define一下
const left =1;
const right = 2;
const down = 3;
const up = 4;

# Called when the node enters the scene tree for the first time.
var gap_y :int = 144; #这个就是最合适的距离和边框的距离。
var gap_x :int = 152;

var locate_status : int = 0;


func _ready() -> void:		
	
	process_noise();
	

func _process(delta: float) -> void:
#速度需要发生变化。上下的浮动同时也要发生变化。那自然而然的角度应该也有变化，这就说明了直接用角度去判断撞墙掉头是非常麻烦的。会出很多bug。
#已经确定了隐形的直线，和一开始说的那样，那就用一个数组去处理状态，直接用int吗？1234四个状态，你用数组本身也要去读取吧？可是数组的上限更高吧，那还不如直接就用一个整数吧，我说实在不行可以当作指针。
	process_move(delta);
		
func change_rad_position():##转角处的变化
	match locate_status :
		down:
			locate_status = left;
			rotation_degrees = 90;
			position = Vector2(0+gap_x,mapHeight+size.x/2)
		left:
			locate_status = up;
			rotation_degrees = 180;
			position = Vector2(0-size.x/2,0+gap_y)
		up:
			locate_status = right ; 
			rotation_degrees = 270;
			position = Vector2(mapWidth-gap_x,0-size.y/2)
		right:
			locate_status = down;
			rotation_degrees = 0;
			position = Vector2(mapWidth+size.x/2,mapHeight - gap_y)

		
func decide_change_rad():#判断改变方向
	if position.x <= 0 - size.x/2 and locate_status == down :
		change_rad_position();
	elif position.y <= 0 - size.x/2 and locate_status == left :
		change_rad_position();
	elif position.x >= mapWidth + size.x/2 and locate_status == up :
		change_rad_position();
	elif position.y >= mapHeight + size.x/2   and locate_status == right :
		change_rad_position();
		
func process_noise(): #噪音的处理
	noise.seed = randi();
	noise.frequency = 0.;
	
func process_move(delta):#运动处理、
	
	
	time += delta;
	var n = noise.get_noise_1d(time)
	match locate_status :
		down:
			position.x -= speed;
			position.y += n * amplitude;
		left:

			position.y -= speed;
			position.x += n*amplitude;
		up:

			position.x += speed;
			position.y += n*amplitude;
		right:

			position.y += speed;
			position.x -= n*amplitude;
	
func 	
	
	
