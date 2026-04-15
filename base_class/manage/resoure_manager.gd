extends Node

## 导出变量
@export var card_scene: PackedScene = preload("res://base_class/card/card.tscn")
@export var spawn_container: Control
@export_dir var base_path: String = "res://game_data/card/" # 资源存放的基础路径

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
	
