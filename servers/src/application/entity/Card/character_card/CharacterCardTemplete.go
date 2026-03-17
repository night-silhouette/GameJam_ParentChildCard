package character_card

import _ "embed"

type CharacterCardTemplate struct {
	ID   int            `json:"id"`
	Info map[string]any `json:"-"`
}

func (c *CharacterCardTemplate) GetID() int {
	return c.ID
}

func (c *CharacterCardTemplate) GetInfo() map[string]any {
	return c.Info
}
