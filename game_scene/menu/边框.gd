extends Control
@export var _animationpalyer : AnimationPlayer;


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	_animationpalyer.play("边框进入")


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass
