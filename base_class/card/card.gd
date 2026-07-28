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




## 延迟 n 秒后发送 SignalBus.time_end




## 将当前卡牌数据导入到卡牌具体页（所有子类通用）
func populate_detail_page(detail_page: Control) -> void:
	if not detail_page:
		return
	
	var name_label = detail_page.get_node_or_null("name") as Label
	var texture_rect = detail_page.get_node_or_null("卡牌纹理") as TextureRect
	var parent_icon = detail_page.get_node_or_null("母牌") as TextureRect
	var child_icon = detail_page.get_node_or_null("子牌") as TextureRect
	var hp_label = detail_page.get_node_or_null("信息标血量/hp") as Label
	var attack_label = detail_page.get_node_or_null("信息标攻击/attack") as Label
	var cost_label = detail_page.get_node_or_null("点ui/energy") as Label
	var skill_effect_label = detail_page.get_node_or_null("信息标技能/技能效果") as Label
	var summon_condition_label = detail_page.get_node_or_null("信息标技能/召唤条件") as Label
	var summon_effect_label = detail_page.get_node_or_null("信息标技能/召唤效果") as Label
	
	# 卡牌名
	if name_label: name_label.text = card_name
	# 卡面纹理
	if texture_rect and card_texture: texture_rect.texture = card_texture
	# 母牌/子牌图标
	var is_parent = not is_sub_card
	if parent_icon: parent_icon.visible = is_parent
	if child_icon: child_icon.visible = not is_parent
	# HP / Attack / 能量消耗
	if hp_label: hp_label.text = str(max_health)
	if attack_label: attack_label.text = str(card_damage)
	if cost_label: cost_label.text = str(skill_charge)
	# 技能效果/召唤条件/召唤效果
	if skill_effect_label:
		skill_effect_label.text = "技能效果：" + skill_description
		skill_effect_label.visible = true
	if is_parent:
		if summon_condition_label: summon_condition_label.visible = false
		if summon_effect_label: summon_effect_label.visible = false
	else:
		if summon_condition_label:
			summon_condition_label.text = "召唤条件：" + notes
			summon_condition_label.visible = true
		if summon_effect_label:
			summon_effect_label.text = "召唤效果：" + sub_card_trigger_effect
			summon_effect_label.visible = true
