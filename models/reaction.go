package models

import (
	"github.com/gofrs/uuid"
)

type Reaction struct {
	Emoji   string      `json:"emoji"`
	UserIDs []uuid.UUID `json:"user_ids"`

	Users               []User `json:"users"`
	IsMadeByCurrentUser bool   `json:"-"`
}
