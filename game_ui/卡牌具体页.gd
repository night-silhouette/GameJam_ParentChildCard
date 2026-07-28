extends Control

@onready var name_label: Label = $"name"
@onready var card_texture_rect: TextureRect = $"卡牌纹理"
@onready var parent_icon: TextureRect = $"母牌"
@onready var child_icon: TextureRect = $"子牌"
@onready var hp_label: Label = $"信息标血量/hp"
@onready var attack_label: Label = $"信息标攻击/attack"
@onready var skill_effect_label: Label = $"信息标技能/技能效果"
@onready var summon_condition_label: Label = $"信息标技能/召唤条件"
@onready var summon_effect_label: Label = $"信息标技能/召唤效果"
@onready var cost_label: Label = $"点ui/energy"


## 导入卡牌数据到详情页（总函数）
## @param card_data: CardResource - 本地卡牌资源数据
func import_card_data(card_data: CardResource) -> void:
	if card_data == null:
		return
	
	# 名称
	name_label.text = card_data.name
	
	# 卡牌纹理
	if card_data.card_texture:
		card_texture_rect.texture = card_data.card_texture
	
	# 母牌/子牌图标
	var is_parent = not card_data.is_sub_card
	parent_icon.visible = is_parent
	child_icon.visible = not is_parent
	
	# HP
	hp_label.text = str(card_data.max_health)
	
	# Attack
	attack_label.text = str(card_data.damage)
	
	# 能量消耗
	cost_label.text = str(card_data.value)
	
	# 技能效果 / 召唤条件 / 召唤效果（仅母牌显示）
	if is_parent:
		skill_effect_label.text = "技能效果：" + card_data.skill_description
		summon_condition_label.text = "召唤条件：" + card_data.notes
		summon_effect_label.text = "召唤效果：" + card_data.sub_card_trigger_effect
		skill_effect_label.visible = true
		summon_condition_label.visible = true
		summon_effect_label.visible = true
	else:
		skill_effect_label.visible = false
		summon_condition_label.visible = false
		summon_effect_label.visible = false
