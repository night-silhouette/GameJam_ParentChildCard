package CardAbstract

type Character interface {
	Card
	Attack(TargetId int)
	Hurt(AttackTempId int, HurtHp float64) //id是谁打的他
	Skill(TargetId int)
	Death(AttackTempId int)
}
