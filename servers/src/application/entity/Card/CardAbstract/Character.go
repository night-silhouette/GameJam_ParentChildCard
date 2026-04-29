package CardAbstract


type Character interface {
	Card
	Attack(tempId int)
	Hurt(tempId int)
	Skill(tempId int)
}
