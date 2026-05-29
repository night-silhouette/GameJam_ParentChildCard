extends TextureRect

@onready var label: Label = $Label

# 使用 set 关键字拦截变量赋值操作
var txt: String = "":
	set(value):
		txt = value
		# 确保 label 节点已经加载完成，防止游戏刚启动时报错
		if is_node_ready() and label:
			label.text = txt

func _ready() -> void:
	# 初始化时先显示一次默认的 txt 内容
	if label:
		label.text = txt
