

extends TextureRect
var x = 576;
@export var Duichen : TextureRect ;
@export var animationpalyer : AnimationPlayer ;

func _ready() -> void:
	#handle_animation();
	pass;

func _process(delta: float) -> void:
	
	handle_movement();
	

func handle_movement() -> void:#大小一致和对游戏界面最中心的中心对称。
	Duichen.position.x=2*x -  position.x -size.x * scale.x;
	Duichen.size.x = size.x*scale.x;
	
	Duichen.position.y = position.y;
	Duichen.size.y = size.y;

func handle_animation() -> void:#单纯播放动画。
	
	animationpalyer.play("边框进入");
	
