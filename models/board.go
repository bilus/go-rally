package models

import (
	"encoding/json"
	"time"

	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gobuffalo/validate/v3/validators"
	"github.com/gofrs/uuid"
)

// Board is used by pop to map your boards database table to your go code.
type Board struct {
	ID                uuid.UUID    `json:"id" db:"id"`
	Name              string       `json:"name" db:"name"`
	Description       nulls.String `json:"description" db:"description"`
	VotingStrategy    `json:"-" db:"-"`
	VotingStrategyRaw json.RawMessage `json:"-" db:"voting_strategy" form:"-"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
	Private           bool            `json:"private" db:"is_private"`

	// User-context fields.
	UserStarred bool `json:"-" db:"-" form:"-"`
	UserIsOwner bool `json:"-" db:"-" form:"-"`

	Posts []Post `json:"posts" db:"-" form:"-"`
}

type VotingStrategy struct {
	BoardMax nulls.Int `json:"board_max"`
}

func DefaultBoard() *Board {
	return &Board{
		VotingStrategy: VotingStrategy{
			BoardMax: nulls.NewInt(10),
		},
	}
}

// String is not required by pop and may be deleted
func (b Board) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Boards is not required by pop and may be deleted
type Boards []Board

func (bs Boards) Slice() []Board {
	return ([]Board)(bs)
}

func (bs *Boards) AddUserContext(user *User) {
	slice := bs.Slice()
	for i, b := range slice {
		slice[i].UserStarred = user.IsBoardStarred(&b)
	}
}

// String is not required by pop and may be deleted
func (b Boards) String() string {
	jb, _ := json.Marshal(b)
	return string(jb)
}

// Validate gets run every time you call a "pop.Validate*" (pop.ValidateAndSave,
// pop.ValidateAndCreate, pop.ValidateAndUpdate) method. This method is not
// required and may be deleted.
func (b *Board) Validate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.Validate(
		&validators.StringIsPresent{Field: b.Name, Name: "Name"},
	), nil
}

// ValidateCreate gets run every time you call "pop.ValidateAndCreate" method.
// This method is not required and may be deleted.
func (b *Board) ValidateCreate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// ValidateUpdate gets run every time you call "pop.ValidateAndUpdate" method.
// This method is not required and may be deleted.
func (b *Board) ValidateUpdate(tx *pop.Connection) (*validate.Errors, error) {
	return validate.NewErrors(), nil
}

// BeforeSave callback will be called before a record is either created or
// updated in the database.
func (b *Board) BeforeSave(*pop.Connection) error {
	data, err := json.Marshal(b.VotingStrategy)
	if err != nil {
		return err
	}
	b.VotingStrategyRaw = data
	return nil
}

// AfterFind callback will be called after a record, or records, has been
// retrieved from the database.
func (b *Board) AfterFind(*pop.Connection) error {
	return json.Unmarshal(b.VotingStrategyRaw, &b.VotingStrategy)
}
