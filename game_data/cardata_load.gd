
extends Node

# 在检查器中点击这个开关开始执行
@export var start_import: bool = false:
	set(value):
		if value: 
			import_csv()
			start_import = false

# 配置路径（请根据你的实际项目目录修改）
const ART_BASE_PATH = "res://素材/使用素材/卡牌/" # 美术素材根目录
const SAVE_DIR = "res://game_data/card/" # 生成的资源存放目录
const CSV_PATH = "res://game_data/card_data/card.csv"        # CSV文件路径

func import_csv():
	# 1. 确保保存目录存在
	if not DirAccess.dir_exists_absolute(SAVE_DIR):
		DirAccess.make_dir_recursive_absolute(SAVE_DIR)

	# 2. 打开 CSV 文件
	var file = FileAccess.open(CSV_PATH, FileAccess.READ)
	if not file:
		printerr("错误：无法打开 CSV 文件，请检查路径：" + CSV_PATH)
		return

	file.get_csv_line() # 跳过第一行表头
	
	var count = 0
	while !file.eof_reached():
		var line = file.get_csv_line()
		if line.size() < 2: continue # 略过空行

		var res_id = int(line[1])
		var file_name = "card_%03d.tres" % res_id
		var full_save_path = SAVE_DIR + file_name
		
		# 3. 加载或创建资源对象
		var res: CardResource
		if FileAccess.file_exists(full_save_path):
			res = load(full_save_path) # 如果已存在则加载，保留手动修改的内容
		else:
			res = CardResource.new()

		# 4. 【核心步骤】调用递归函数查找纹理
		var target_image_name = str(res_id) + ".png"
		var actual_img_path = _find_file_recursive(ART_BASE_PATH, target_image_name)
		
		if actual_img_path != "":
			res.card_texture = load(actual_img_path)
		else:
			# 如果没找到，控制台会打印黄色警告，方便你检查哪个ID缺图
			push_warning("ID %d: 递归查找失败，未在 %s 下找到 %s" % [res_id, ART_BASE_PATH, target_image_name])

		# 5. 同步 CSV 中的其他数据
		res.name = line[0]
		res.id = res_id
		res.damage = int(line[2]) if line[2] != "" else 0
		res.max_health = int(line[3]) if line[3] != "" else 0
		res.is_sub_card = (line[4] == "子")
		res.is_combat_card = (line[5] == "战")
		res.effect_description = line[6]
		res.notes = line[7]

		# 6. 保存资源
		ResourceSaver.save(res, full_save_path)
		count += 1
	
	print("成功处理 %d 张卡牌资源！" % count)

# --- 递归查找工具函数 ---
func _find_file_recursive(dir_path: String, target_name: String) -> String:
	var dir = DirAccess.open(dir_path)
	if not dir:
		return ""
	
	dir.list_dir_begin()
	var file_or_dir = dir.get_next()
	
	while file_or_dir != "":
		# 忽略隐藏文件（如 .git, .import 等）
		if file_or_dir.begins_with("."):
			file_or_dir = dir.get_next()
			continue
			
		var full_path = dir_path.path_join(file_or_dir)
		
		if dir.current_is_dir():
			# 如果当前是文件夹，递归调用自身进入下一层
			var found = _find_file_recursive(full_path, target_name)
			if found != "": 
				return found # 如果在子文件夹里找到了，直接向上层返回路径
		else:
			# 如果当前是文件，比对文件名
			if file_or_dir == target_name:
				return full_path
				
		file_or_dir = dir.get_next()
	
	return "" # 全找遍了也没找到，返回空
