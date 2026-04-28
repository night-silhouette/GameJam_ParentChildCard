package CardAbstract

type Character interface {
	Card
	Attack()
	Hurt()
	Skill()
}
