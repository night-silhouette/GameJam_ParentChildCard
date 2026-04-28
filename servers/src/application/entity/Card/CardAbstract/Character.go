package CardAbstract

import "pcc_card/application/entity/BattleData"

type Character interface {
	Card
	Attack(w BattleData.Where)
	Hurt()
	Skill()
}
