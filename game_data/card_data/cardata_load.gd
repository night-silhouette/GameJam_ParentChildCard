# CardImporter.gd
@tool
extends Node

@export_group("Paths")
@export_file("*.csv") var csv_path: String = "res://data/Card.csv"
@export_dir var output_dir: String = "res://resources/cards/"
@export_dir var texture_dir: String = "res://assets/cards/"
@export_file("*.png", "*.jpg") var default_texture_path: String = "res://assets/cards/default_texture.png"

@export_group("Controls")
## 点击勾选此项立即开始导入（文件名绝对唯一，无等级后缀）
@export var trigger_import: bool = false:
	set(value):
		if value == true:
			_start_import()
			trigger_import = false 

func _start_import() -> void:
	# print("========= 开始精准导入卡牌数据 (干净 ID 命名模式) =========")
	
	# 确保输出目录存在
	if not DirAccess.dir_exists_absolute(output_dir):
		DirAccess.make_dir_recursive_absolute(output_dir)
		
	# 加载默认纹理保底
	var default_texture: Texture2D = null
	if ResourceLoader.exists(default_texture_path):
		default_texture = load(default_texture_path)

	var file := FileAccess.open(csv_path, FileAccess.READ)
	if not file:
		push_error("无法打开 CSV 文件: " + csv_path)
		return
		
	# 存储列名到索引的字典映射
	var header_map = {}
	var imported_count = 0

	# 辅助取值匿名函数：根据表头名字获取行里对应的数据
	var get_val = func(row: PackedStringArray, column_name: String, default_str: String = "") -> String:
		for key in header_map.keys():
			if column_name in key:
				var idx = header_map[key]
				if idx < row.size():
					return row[idx].strip_edges()
		return default_str

	# 逐行读取 CSV
	while !file.eof_reached():
		var row = file.get_csv_line()
		
		# 1. 过滤完全没有数据的空行
		if row.size() == 0 or (row.size() == 1 and row[0].strip_edges() == ""):
			continue
			
		var first_cell = row[0].strip_edges()
		
		# 2. 严格跳过各种格式的分界线行
		if "分界线" in first_cell or "——" in first_cell or (row.size() > 2 and "——" in row[2]):
			continue

		# 3. 如果遇到包含“名称”或“ID”的行，说明这是表头行，刷新或初始化列索引
		if "名称" in first_cell or "ID" in row:
			header_map.clear() # 清空旧的，写入当前最新部分的列顺序
			for i in range(row.size()):
				var h_name = row[i].strip_edges()
				if h_name != "":
					header_map[h_name] = i
			# print("[表头刷新] 检测到列顺序更新: ", header_map)
			continue # 表头行不需要当作卡牌导入，跳过
			
		# 如果还没有读取到任何有效的表头，则无法解析，继续往下读
		if header_map.is_empty():
			continue

		# 4. 提取卡牌核心字段
		var card_name = get_val.call(row, "名称")
		var card_id_str = get_val.call(row, "ID")
		
		# 如果名字和 ID 同时为空，说明是噪音空行，跳过
		if card_name == "" and card_id_str == "":
			continue
			
		var card_id = card_id_str.to_int()
		
		# 提取其他字段
		var card_level_str = get_val.call(row, "等级")
		var card_level = 1 if card_level_str == "" else card_level_str.to_int()
		
		var img_name = get_val.call(row, "图片")
		var damage = get_val.call(row, "伤害").to_int()
		var init_hp = get_val.call(row, "初始血量")
		if init_hp == "": init_hp = "0"
		var max_hp = get_val.call(row, "初始生命上限")
		if max_hp == "": max_hp = "0"
		
		var skill_charge = get_val.call(row, "skillCharge").to_int()
		var skill_use_num = get_val.call(row, "skillcardUseNum").to_int()
		var value_num = get_val.call(row, "价值").to_int()
		
		var notes = get_val.call(row, "注释")
		var skill_desc = get_val.call(row, "技能描述")
		
		# 5. 精准转换布尔值类型：如果文本包含"战斗牌"/"子牌"则设为 true，否则为 false
		var is_combat = "战斗牌" in get_val.call(row, "战斗牌/法术牌")
		var is_sub = "子牌" in get_val.call(row, "子牌/母牌")
		var sub_effect = get_val.call(row, "子牌归属触发效果")

		# 6. 【纯净命名】文件命名严格采用你要求的干净格式：card_ID.tres
		var file_name = "card_" + str(card_id) + ".tres"
		var res_file_path = output_dir.path_join(file_name)
		
		# print("正在导出 -> 文件: %s | 名称: %s | 等级: %d | 战斗牌: %s | 子牌: %s" % [file_name, card_name, card_level, str(is_combat), str(is_sub)])

		var card_res: CardResource
		if ResourceLoader.exists(res_file_path):
			card_res = load(res_file_path)
		else:
			card_res = CardResource.new()
			
		# 数据赋值给 CardResource
		card_res.name = card_name
		card_res.id = card_id
		card_res.level = card_level
		card_res.value = value_num
		card_res.damage = damage
		card_res.initial_health = init_hp.to_int()
		card_res.max_health = max_hp.to_int()
		card_res.skill_charge = skill_charge
		card_res.skill_card_use_num = skill_use_num
		card_res.is_combat_card = is_combat
		card_res.is_sub_card = is_sub
		card_res.notes = notes
		card_res.skill_description = skill_desc
		card_res.sub_card_trigger_effect = sub_effect
		card_res.texture_filename = img_name
		
		# 7. 图片匹配与默认保底逻辑
		var texture_loaded = false
		if img_name != "":
			var full_img_path = texture_dir.path_join(img_name)
			if ResourceLoader.exists(full_img_path):
				card_res.card_texture = load(full_img_path)
				texture_loaded = true
		
		if not texture_loaded:
			var possible_extensions = [".png", ".jpg", ".jpeg", ".tga"]
			for ext in possible_extensions:
				var id_img_path = texture_dir.path_join(str(card_id) + ext)
				if ResourceLoader.exists(id_img_path):
					card_res.card_texture = load(id_img_path)
					texture_loaded = true
					break
					
		if not texture_loaded:
			card_res.card_texture = default_texture
				
		# 写入磁盘
		var save_result = ResourceSaver.save(card_res, res_file_path)
		if save_result == OK:
			imported_count += 1
		else:
			push_error("无法保存文件: " + res_file_path)
			
	file.close()
	# print("========= 导入完成！所有卡牌已成功直接存入 %s 文件夹，共生成 %d 个资源 =========" % [output_dir, imported_count])
