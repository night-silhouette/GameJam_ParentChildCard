class_name card
extends Control

var card_name: String = ""
var card_texture: Texture2D
var id: int = 0
var level: int = 1
var value: int = 0

var card_damage: int = 0
var initial_health: int = 0
var max_health: int = 0

var skill_charge: int = 0
var skill_card_use_num: int = 0

var is_combat_card: bool = false
var is_sub_card: bool = false

var skill_description: String = ""
var notes: String = ""
var sub_card_trigger_effect: String = ""
