# CardResource.gd
class_name CardResource extends Resource

@export_group("Visuals")
# 在这里添加纹理字段，你可以在编辑器里手动拖入 .png 文件
@export var card_texture: Texture2D 

@export_group("Data from CSV")
@export var name: String
@export var id: int
@export var damage: int
@export var max_health: int
@export var is_sub_card: bool
@export var is_combat_card: bool

@export_group("Descriptions")
@export_multiline var effect_description: String
@export_multiline var notes: String
	
