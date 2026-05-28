extends Node

# 本地背包数据
var card_list: Array = []

@export var bag_zone = Global.ZONE_CARD.BAG_ZONE
@export var sell_zone = Global.ZONE_CARD.SELL_ZONE


func _ready() -> void:
	SignalBus.get_card_bag.connect(_get_card_bag)


func _get_card_bag(card_list):

	card_list.clear()
	for item in card_list:

		var data = {
			"stuff_id": item["stuff_id"],
			"card_id": item["card_id"],
			"price": item["price"],
			"zone": bag_zone
		}

		card_list.append(data)

# 根据 stuff_id 修改 zone 为 sell_zone
func move_to_sell_zone(stuff_id: int) -> bool:

	for item in card_list:

		if item["stuff_id"] == stuff_id:
			item["zone"] = sell_zone

			return true

	return false
