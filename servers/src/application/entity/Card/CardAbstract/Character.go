package CardAbstract

import "pcc_card/application/entity/BattleData"

type Character interface {
	Card
	BtCry()

	Attack(TargetId int) bool
	Hurt(AttackTempId int, HurtHp float64, category BattleData.ValueChange) //id是谁打的他
	NoSourceHurt(HurtHp float64, category BattleData.ValueChange)
	Skill(TargetId int) bool
	Death(AttackTempId int) //如果,杀死者是无主的,就传入-1

}
