extends Control

var zones: Array[Control]   # 所有可进入区域
@onready var display: TextureRect = $"卡牌纹理"
var temp_id ;
var hp ;
var damage  ;
var buff_id ;
var zone ;

var dragging = false
var click_offset = Vector2.ZERO
var original_position = Vector2.ZERO
var current_zone: Control = null

func _ready():
	original_position = global_position
func update_card(base_res) -> void:
	display.texture = base_res.texture
	if base_res.is_combat_card :
		hp = base_res.hp;
		damage = base_res.damage;
	temp_id = base_res.temp_id;
	buff_id = base_res.buff_id;
func update():
	pass;
	
func free_card():
	display.texture = null;
	hp = null;
	damage = null;
	temp_id = null;
	buff_id = null;


func _gui_input(event):
	if event is InputEventMouseButton:
		if event.button_index == MOUSE_BUTTON_LEFT:
			if event.pressed:
				# 开始拖拽
				dragging = true
				# 记录鼠标相对于节点左上角的偏移，避免卡牌中心猛地跳到鼠标位置
				click_offset = get_global_mouse_position() - global_position
				# 拖拽时置顶显示
				move_to_front() 
			else:
				# 停止拖拽
				dragging = false
				_on_card_dropped()

func _process(_delta):
	if dragging:
		# 更新位置：当前鼠标全局位置减去初始点击偏移
		global_position = get_global_mouse_position() - click_offset
		check_zone();
		
func check_zone():
	var center = get_global_rect().get_center()
	current_zone = null

	for z in zones:
		if z.get_global_rect().has_point(center):
			current_zone = z
			return
func _on_card_dropped():
	if current_zone:
		zone = current_zone.name   # 或者你自己定义
		_snap_to_zone(current_zone)
	else:
		_back()

func _snap_to_zone(target: Control):
	var tween = create_tween()
	tween.tween_property(
		self,
		"global_position",
		target.global_position,
		0.2
	).set_trans(Tween.TRANS_QUART).set_ease(Tween.EASE_OUT)
func _back():

	var tween = create_tween()
	tween.tween_property(self, "global_position", original_position, 0.2).set_trans(Tween.TRANS_QUART).set_ease(Tween.EASE_OUT)
