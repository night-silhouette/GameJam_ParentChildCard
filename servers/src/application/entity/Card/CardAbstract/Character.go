package CardAbstract

import "pcc_card/application/entity/BattleData"

type Character interface {
	Card
	BtCry()

	Attack(TargetId int)
	Hurt(AttackTempId int, HurtHp float64, category BattleData.ValueChange) //id是谁打的他
	Skill(TargetId int) bool
	Death(AttackTempId int)
	RoundEnd()
}
