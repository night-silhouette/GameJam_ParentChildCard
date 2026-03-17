package skill_card

type SkillCardTemplate struct {
	ID   int            `json:"id"`
	Info map[string]any `json:"-"`
}

func (c *SkillCardTemplate) SetInfo(info map[string]any) {
	c.Info = info
}
func (c *SkillCardTemplate) GetInfo() map[string]any {
	return c.Info
}
