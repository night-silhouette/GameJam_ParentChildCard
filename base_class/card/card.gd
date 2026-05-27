class_name card
extends Control

@onready var display: TextureRect = $"卡牌纹理"
var temp_id ;
var hp ;
var damage  ;
var buff_id ;
var zone ;
var spell_des ;

func update_card_data(base_res) -> void:
	display.texture = base_res.texture
	if base_res.is_combat_card :
		hp = base_res.hp;
		damage = base_res.damage;
	temp_id = base_res.temp_id;
	buff_id = base_res.buff_id;
	zone = base_res.zone;
	spell_des = base_res.spell_des;

func free_card():
	display.texture = null;
	hp = null;
	damage = null;
	temp_id = null;
	buff_id = null;
	
