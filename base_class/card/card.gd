extends Control

@onready var display: TextureRect = $"卡牌纹理"
var net_ID ;
var hp ;
var atk  ;

func init_card(base_res: CardResource) -> void:
	display.texture = base_res.card_texture
	if base_res.is_combat_card :
		hp = base_res.max_health;
		atk = base_res.damage;
func update():
	pass;
	
func free_card():
	display.texture = null;
	hp = null;
	atk = null;
	net_ID = null;
