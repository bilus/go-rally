package models

import (
	"encoding/json"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	"time"
)
// BoardMember is used by pop to map your board_members database table to your go code.
type BoardMember struct {
    ID uuid.UUID `json:"id" db:"id"`
    BoardID uuid.UUID `json:"board_id" db:"board_id"`
    UserID uuid.UUID `json:"user_id" db:"user_id"`
    IsOwner bool `json:"is_owner" db:"is_owner"`
    CreatedAt time.Time `json:"created_at" db:"created_at"`
    UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (b BoardMember) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// BoardMembers is not required by pop and may be deleted
type BoardMembers []BoardMember

// String is not required by pop and may be deleted
func (b BoardMembers) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (b *BoardMember) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *BoardMember) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *BoardMember) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
