extends Node2D

var velocity : Vector2
var gravity := 600
var rotate_speed := 5


func _ready():

	# 随机方向
	velocity = Vector2(randf_range(-120,120), randf_range(-250,-150))


func _process(delta):

	velocity.y += gravity * delta

	position += velocity * delta

	rotation += rotate_speed * delta
	
func _on_timer_timeout():
	queue_free()
