package models

import "github.com/gofrs/uuid"

type Reaction struct {
	Key   string
	Count int
	Users []uuid.UUID
}
