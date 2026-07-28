extends Control
var price: int 
var card_id: int
var goods_id: int
var discount : float


# Called when the node enters the scene tree for the first time.
func _ready() -> void:
	pass # Replace with function body.


# Called every frame. 'delta' is the elapsed time since the previous frame.
func _process(delta: float) -> void:
	pass

func set_up(data:Dictionary):
	price = data.get("price")
	card_id = data.get("card_id")
	goods_id = data.get("goods_id")
	discount = data.get("discount")
	$coin/num.text = str(price)
	if discount == 1:
		$"折扣".visible = false
	else:
		$"折扣".visible = true
		$"折扣/打折".text = str(roundi((discount - 1) * 100)) + "%"
	var res = InventoryManager._find_card_resource_by_id(card_id)
	data["resource"]  = res
	$shopcard.setup(data)
	
	
