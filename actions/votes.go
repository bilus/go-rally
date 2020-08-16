package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
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

	n, err := tx.RawQuery("UPDATE users SET votes = votes - 1 WHERE votes > 0 AND id = ?", u.ID).ExecWithCount()
	if err != nil {
		return err
	}
	if n == 0 {
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("votes/fail.js"))
	}

	// Allocate an empty Post
	post := &models.Post{}

	// To find the Post the parameter post_id is used.
	if err := tx.Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	post.Votes++ // TODO: Race condition.

	vote := &models.Vote{PostID: post.ID, UserID: u.ID}
	verrs, err := tx.ValidateAndCreate(vote)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		c.Set("errors", verrs)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("error.js"))
	}

	if err := tx.UpdateColumns(post, "votes", "updated_at"); err != nil {
		return fmt.Errorf("upvoting failed: %v", err)
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

	var vote models.Vote
	err = tx.Where("post_id = ? AND user_id = ?", c.Param("post_id"), u.ID).First(&vote)
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

	// Allocate an empty Post
	post := &models.Post{}

	// To find the Post the parameter post_id is used.
	if err := tx.Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	post.Votes-- // TODO: Race condition.

	if err := tx.UpdateColumns(post, "votes", "updated_at"); err != nil {
		return fmt.Errorf("upvoting failed: %v", err)
	}

	err = tx.RawQuery("UPDATE users SET votes = votes + 1 WHERE id = ?", u.ID).Exec()
	if err != nil {
		return err
	}

	c.Set("post", post)
	return c.Render(http.StatusOK, r.JavaScript("votes/create.js"))
}
