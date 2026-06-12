extends LineEdit


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	self.add_theme_color_override(
	"font_color",
	Color.DARK_KHAKI
)
#加入数字判定。
