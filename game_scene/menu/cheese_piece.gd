extends Node2D

var velocity : Vector2
var gravity := 1500  # 增加重力感
var rotate_speed : float

func _ready():
	# 随机左右方向并赋予较强的横向力
	var side = 1 if randf() > 0.5 else -1
	velocity.x = randf_range(200, 400) * side
	velocity.y = randf_range(-400, -200) # 向上跳跃的力
	
	# 随机旋转方向和速度
	rotate_speed = randf_range(-8.0, 8.0)

func _process(delta):
	# 标准物理公式：速度受重力影响
	velocity.y += gravity * delta
	position += velocity * delta
	rotation += rotate_speed * delta

func _on_timer_timeout():
	queue_free()
