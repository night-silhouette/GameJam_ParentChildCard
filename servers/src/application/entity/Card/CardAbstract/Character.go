package CardAbstract

type Character interface {
	Card
	Attack(tempId int)
	Hurt(tempId int, HurtHp int) //id是谁打的他
	Skill(tempId int)
	Death(tempId int)
}
