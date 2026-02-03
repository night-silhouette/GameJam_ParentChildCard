@tool
extends TextureRect
var x = 576;
@export var Duichen : TextureRect ;


func _process(delta: float) -> void:
	

	
	Duichen.position.x=2*x -  position.x - size.x;
	Duichen.size.x = size.x;
	
	Duichen.position.y = position.y;
	Duichen.size.y = size.y;
	
