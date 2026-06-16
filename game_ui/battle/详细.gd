extends Control

func update_card_data(card):
	$hp.text = str(card.get("hp", ""))
	$damage.text = str(card.get("damage", ""))
	
	var res = card.get("resouce")
	if res:
		$"详细".text = str(res.get("skill_description") if res.get("skill_description") else res.get("notes", ""))
	
