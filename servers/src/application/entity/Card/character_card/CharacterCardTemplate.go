package character_card

import _ "embed"

type CharacterCardTemplate struct {
	ID   int            `json:"id"`
	Info map[string]any `json:"-"`
}

func (c *CharacterCardTemplate) SetInfo(info map[string]any) {
	c.Info = info
}
func (c *CharacterCardTemplate) GetInfo() map[string]any {
	return c.Info
}
