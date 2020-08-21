package models

import (
	"encoding/json"
	"io"
	"os"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Attachment is used by pop to map your attachments database table to your go code.
type Attachment struct {
	ID        uuid.UUID `json:"id" db:"id"`
	PostID    uuid.UUID `json:"post_id" db:"post_id"`
	Filename  string    `json:"filename" db:"filename"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// String is not required by pop and may be deleted
func (a Attachment) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Attachments is not required by pop and may be deleted
type Attachments []Attachment

// String is not required by pop and may be deleted
func (a Attachments) String() string {
	ja, _ := json.Marshal(a)
	return string(ja)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (a *Attachment) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: a.Filename, Name: "Filename"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (a *Attachment) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (a *Attachment) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// Open opens the attachment for reading.
// Note: The returned ReadCloser must be closed.
func (a *Attachment) Open() (io.ReadCloser, error) {
	return os.Open(a.ID.String())
}

func (a *Attachment) Save(r io.Reader) error {
	w, err := os.Create(a.ID.String())
	if err != nil {
		return err
	}
	_, err = io.Copy(w, r)
	return err
}
