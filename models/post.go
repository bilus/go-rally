package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Post is used by pop to map your posts database table to your go code.
type Post struct {
	ID    uuid.UUID `json:"id" db:"id"`
	Title string    `json:"title" db:"title"`
	Body  string    `json:"body" db:"body"`

	// Votes cache values in Redis.
	Votes int `json:"votes" db:"votes" form:"-"`

	// CommentCount caches post comment count.
	CommentCount int `json:"comment_count" db:"comment_count" form:"-"`

	CreatedAt time.Time `json:"created_at" db:"created_at" form:"-"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at" form:"-"`

	AuthorID uuid.UUID `json:"-" db:"author_id"`
	Author   User      `json:"author,omitempty" belongs_to:"user" db:"-"`

	Draft     bool `json:"draft" db:"draft"`
	Archived  bool `json:"archived" db:"archived"`
	Anonymous bool `json:"anonymous" db:"anonymous"` // TODO: When set, Author should not be marshalled to JSON.

	BoardID uuid.UUID `json:"-" db:"board_id"`

	// Optionally looaded.
	Reactions []Reaction `json:"reactions,omitempty" db:"-" from:"-"`
}

// String is not required by pop and may be deleted
func (p Post) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Posts is not required by pop and may be deleted
type Posts []Post

// String is not required by pop and may be deleted
func (p Posts) String() string {
	jp, _ := json.Marshal(p)
	return string(jp)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave, pop.ValidateAndCreate, pop.ValidateAndUpdate) method.
// This method is not required and may be deleted.
func (p *Post) Validate(tx *pop.Connection) (*validate.Errors, error) {
	if p.Draft {
		return validate.Validate(), nil
	}

	return validate.Validate(
		&validators.StringIsPresent{Field: p.Title, Name: "Title"},
		&validators.StringIsPresent{Field: p.Body, Name: "Body"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (p *Post) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (p *Post) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}
