extends Control

func update_card_data(card):
	$hp.text = str(card.get("hp", ""))
	$damage.text = str(card.get("damage", ""))
	$"详细".text = str(card.get("spell_des", ""))
