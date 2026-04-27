extends Node
signal show_card;
## 导出变量
@export var card_scene: PackedScene = preload("res://base_class/card/card.tscn")
@export var spawn_container: Control
@export_dir var base_path: String = "res://game_data/card/" # 资源存放的基础路径

var card_list :Array = [];
func _ready() -> void:
	SignalBus.self_cards_updated.connect(_self_cards_updated)
	
## 核心功能：通过 ID 生成resoure
func querry_resoure_by_id(card_id: int) -> Resource:
	# 1. 格式化路径，%03d 会将 1 转换为 001，将 12 转换为 012
	var file_name = "card_%03d.tres" % card_id
	var full_path = base_path.path_join(file_name)
	
	# 2. 检查文件是否存在
	if not FileAccess.file_exists(full_path):
		push_error("CardManager: 找不到资源文件 -> " + full_path)
		return null
	
	# 3. 加载资源
	var data = load(full_path) as CardResource
	if data == null:
		push_error("CardManager: 资源加载失败或类型错误 -> " + full_path)
		return null
	
	return data;
	
func _self_cards_updated(data):
	if card_list.is_empty():
		card_list = _init_cards(data);
	else :
		_update_cards(data);
	show_card.emit();
	
func _init_cards(cards: Array) :
	var result: Array = []
	for card in cards:
		var new_card = card.duplicate()
		var resoure_data = querry_resoure_by_id(int(card.id));
		new_card["is_combat_card"] = resoure_data.is_combat_card;
		new_card["is_sub_card"] = resoure_data.is_sub_card;
		new_card["texture"] = resoure_data.card_texture;
		result.append(new_card)
	return result;
	
func find_card_by_key(cards: Array, target_dict: Dictionary, key_to_match = "temp_id"):
	# 获取我们要找的那个目标值
	var target_value = target_dict.get(key_to_match)
	for card in cards:
		# 匹配键值
		if card.get(key_to_match) == target_value:
			return card # 返回找到的字典引用（内存地址）
			
	return null # 没找到则返回空
	
func _update_cards(data : Array) :
	# 建议倒序遍历，因为在循环中删除数组元素，正序会导致索引错乱
	for i in range(card_list.size() - 1, -1, -1):
		var card = card_list[i]
		var result = find_card_by_key(data, card, "id")
		if result == null:
		# 找不到匹配，从数组中移除
			card_list.remove_at(i) 
		# 此时，如果没有其他变量引用这个 card 字典，它会被 Godot 自动回收
		else:
		# 找到了，进行后续逻辑
			card.hp = result.hp;
			card.damage = result.damage;
			card.buff_id = result.buff_id;
			
				
