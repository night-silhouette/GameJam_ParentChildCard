package character_card

import _ "embed"

type CharacterCardTemplate struct {
	ID   int            `json:"id"`
	Info map[string]any `json:"-"`
}
