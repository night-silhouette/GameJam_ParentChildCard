extends Node2D

@onready var anim_player: AnimationPlayer = $AnimationPlayer
@onready var card_attack: AnimatedSprite2D = $card_attack
@onready var card_spell: AnimatedSprite2D = $card_spell

func _ready() -> void:
	anim_player.play("card_attack")
