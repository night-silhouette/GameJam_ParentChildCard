extends TextureRect

# 修正：改用 @onready 自动获取子节点中的 Label
# 如果你的 Label 节点名字就叫 "Label"，这样写最稳妥
@onready var num: Label = $Label

# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	# 确保安全：如果信号已经连过了就不再重复连接
	if not InventoryManager.gold_updated.is_connected(_gold_updated):
		InventoryManager.gold_updated.connect(_gold_updated)


# 当金币更新时触发
func _gold_updated(data) -> void:
	if num:
		# 使用 str() 把传入的 data 转换为字符串并赋值给 text
		num.text = str(data)
		
