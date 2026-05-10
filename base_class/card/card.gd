extends Control
   # 所有可进入区域
@onready var display: TextureRect = $"卡牌纹理"
var temp_id ;
var hp ;
var damage  ;
var buff_id ;
var current_zone ;

var dragging = false
var click_offset = Vector2.ZERO
var original_position = Vector2.ZERO

signal request_drag(temp_id);


func _ready():
	original_position = global_position
func update_card_data(base_res) -> void:
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
