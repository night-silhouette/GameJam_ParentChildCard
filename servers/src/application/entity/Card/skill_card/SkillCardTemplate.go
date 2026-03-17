package skill_card

type SkillCardTemplate struct {
	ID   int            `json:"id"`
	Info map[string]any `json:"-"`
}
