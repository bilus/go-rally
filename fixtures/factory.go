package fixtures

import (
	"log"
	"rally/models"

	"github.com/Pallinder/go-randomdata"
	"github.com/gobuffalo/nulls"
	"github.com/gobuffalo/pop/v5"
)

type Factory struct {
	db *pop.Connection
}

func NewFactory(db *pop.Connection) Factory {
	return Factory{
		db: db,
	}
}

func (f Factory) ValidPost(board *models.Board, author *models.User) *models.Post {
	return &models.Post{
		Title:    randomdata.SillyName(),
		Body:     randomdata.Paragraph(),
		Votes:    randomdata.Number(1, 100),
		BoardID:  board.ID,
		AuthorID: author.ID,
		Author:   *author,
	}
}

func (f Factory) EmptyPostDraft(board *models.Board, author *models.User) *models.Post {
	p := f.ValidPost(board, author)
	p.Draft = true
	p.Title = ""
	p.Body = ""
	return p
}

func (f Factory) MustCreatePost(post *models.Post) *models.Post {
	verrs, err := f.db.ValidateAndCreate(post)
	if verrs.HasAny() {
		log.Fatalf("post validation errors: %v", verrs.String())
	}
	if err != nil {
		log.Fatalf("error creating post: %v", err)
	}
	return post
}

func (f Factory) MustCreatePosts(posts ...*models.Post) []*models.Post {
	for _, post := range posts {
		f.MustCreatePost(post)
	}
	return posts
}

func (f Factory) ValidPosts(n int, board *models.Board, author *models.User, drafts bool) []*models.Post {
	ps := make([]*models.Post, n)
	for i := range ps {
		ps[i] = f.ValidPost(board, author)
		ps[i].Draft = drafts
	}
	return ps
}

func (f Factory) MustCreateUser() *models.User {
	u := &models.User{
		Email:                randomdata.Email(),
		Password:             "password",
		PasswordConfirmation: "password",
		GoogleUserID:         nulls.NewString("123"),
		Votes:                5,
	}

	verrs, err := u.Create(f.db)
	if verrs.HasAny() {
		log.Fatalf("user validation errors: %v", verrs.String())
	}
	if err != nil {
		log.Fatalf("error creating user: %v", err)
	}

	return u
}

func (f Factory) MustCreateBoard() *models.Board {
	b := &models.Board{
		Name: randomdata.SillyName(),
	}
	verrs, err := f.db.ValidateAndCreate(b)
	if verrs.HasAny() {
		log.Fatalf("board validation errors: %v", verrs.String())
	}
	if err != nil {
		log.Fatalf("error creating board: %v", err)
	}
	return b
}