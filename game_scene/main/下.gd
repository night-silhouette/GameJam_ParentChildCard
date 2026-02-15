
extends TextureRect
var y = 324;
var x = 576;
@export var Duichen : TextureRect ;
@export var animationpalyer : AnimationPlayer ;

func _ready() -> void:
	handle_animation();

func _process(delta: float) -> void:

	handle_movement();
	


func handle_movement() -> void:
	Duichen.position.x=2*x -  position.x;
	Duichen.size.x = size.x;
	
	Duichen.position.y =2*y - position.y;
	Duichen.size.y = size.y;
	
func handle_animation() -> void:
	animationpalyer.play("边框进入");
	

	

	

	

	
