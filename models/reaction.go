package models

import (
	"github.com/gofrs/uuid"
)

type Reaction struct {
	Emoji string      `json:"emoji"`
	Users []uuid.UUID `json:"user_ids"`

	IsMadeByCurrentUser bool `json:"-"`
}
