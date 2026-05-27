# CardResource.gd
class_name CardResource extends Resource

@export_group("Visuals")
@export var card_texture: Texture2D 
@export var texture_filename: String = ""

@export_group("Base Data")
@export var name: String = ""
@export var id: int = 0
@export var level: int = 1
@export var value: int = 0

@export_group("Combat Stats")
@export var damage: int = 0
@export var initial_health: int = 0
@export var max_health: int = 0

@export_group("Action Settings")
@export var skill_charge: int = 0
@export var skill_card_use_num: int = 0

@export_group("Card Types")
@export var is_combat_card: bool = false
@export var is_sub_card: bool = false

@export_group("Descriptions & Effects")
@export_multiline var skill_description: String = ""
@export_multiline var notes: String = ""
@export_multiline var sub_card_trigger_effect: String = ""
