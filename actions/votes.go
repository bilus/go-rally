package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
)

// VotesCreate upvotes a post
func VotesCreate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	u, err := CurrentUser(c)
	if err != nil {
		return err
	}

	postId, err := uuid.FromString(c.Param("post_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	n, err := tx.RawQuery("UPDATE users SET votes = votes - 1 WHERE votes > 0 AND id = ?", u.ID).ExecWithCount()
	if err != nil {
		return err
	}
	RefreshCurrentUser(c) // Number of votes changed.
	if n == 0 {
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("votes/fail.js"))
	}

	err = tx.RawQuery("UPDATE posts SET votes = votes + 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?", postId).Exec()
	if err != nil {
		return fmt.Errorf("upvoting failed: %v", err)
	}

	vote := &models.Vote{PostID: postId, UserID: u.ID}
	verrs, err := tx.ValidateAndCreate(vote)
	if err != nil {
		return fmt.Errorf("upvoting failed: %v", err)
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("error.js"))
	}

	// Allocate an empty Post
	post := &models.Post{}

	// To find the Post the parameter post_id is used.
	if err := tx.Find(post, postId); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("post", post)
	return c.Render(http.StatusOK, r.JavaScript("votes/create.js"))
}

// VotesDestroy downvotes a post
func VotesDestroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	u, err := CurrentUser(c)
	if err != nil {
		return err
	}

	postId, err := uuid.FromString(c.Param("post_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	var vote models.Vote
	err = tx.Where("post_id = ? AND user_id = ?", postId, u.ID).First(&vote)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			return c.Render(http.StatusUnprocessableEntity, r.JavaScript("votes/fail.js"))
		}
		return err
	}

	err = tx.Destroy(&vote)
	if err != nil {
		return err
	}

	err = tx.RawQuery("UPDATE posts SET votes = votes - 1, updated_at = CURRENT_TIMESTAMP WHERE id = ?", postId).Exec()
	if err != nil {
		return fmt.Errorf("downvoting failed: %v", err)
	}

	err = tx.RawQuery("UPDATE users SET votes = votes + 1 WHERE id = ?", u.ID).Exec()
	if err != nil {
		return fmt.Errorf("downvoting failed: %v", err)
	}
	RefreshCurrentUser(c) // Number of votes changed.

	post := &models.Post{}
	if err := tx.Find(post, postId); err != nil {
		return err
	}
	c.Set("post", post)
	return c.Render(http.StatusOK, r.JavaScript("votes/create.js"))
}

// VotesCreate upvotes a post
// func VotesCreate2(c buffalo.Context) error {
// 	// Get the DB connection from the context
// 	tx, ok := c.Value("tx").(*pop.Connection)
// 	if !ok {
// 		return fmt.Errorf("no transaction found")
// 	}

// 	u, err := CurrentUser(c)
// 	if err != nil {
// 		return err
// 	}

// 	postId, err := uuid.FromString(c.Param("post_id"))
// 	if err != nil {
// 		return c.Error(http.StatusNotFound, err)
// 	}

// 	// To find the Post the parameter post_id is used.
// 	post := &models.Post{}
// 	if err := tx.Eager().Find(post, postId); err != nil {
// 		return c.Error(http.StatusNotFound, err)
// 	}

// 	upvoted, err := post.Board.VotingStrategy.Upvote(tx, u, p)
// 	if err != nil {
// 		return err
// 	}
// 	if !upvoted {
// 		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("votes/error.js"))
// 	}

// 	// To find the Post the parameter post_id is used.
// 	if err := tx.Find(post, postId); err != nil {
// 		return c.Error(http.StatusNotFound, err)
// 	}

// 	c.Set("post", post)
// 	return c.Render(http.StatusOK, r.JavaScript("votes/create.js"))
// }
