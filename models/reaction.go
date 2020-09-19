package models

import "github.com/gofrs/uuid"

type Reaction struct {
	Emoji string
	Users []uuid.UUID
}
