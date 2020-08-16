package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

// VotesCreate upvotes a post
func VotesCreate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Post
	post := &models.Post{}

	// To find the Post the parameter post_id is used.
	if err := tx.Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	post.Votes++ // TODO: Race condition.

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

	c.Set("post", post)
	return c.Render(http.StatusOK, r.JavaScript("votes/create.js"))
}
