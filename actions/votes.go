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

	// Allocate an empty Post
	post := &models.Post{}

	// To find the Post the parameter post_id is used.
	if err := tx.Eager().Find(post, postId); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	postVotes, upvoted, err := post.Board.Upvote(models.Redis, u, post)
	if err != nil {
		return err
	}

	if upvoted {
		post.Votes = postVotes
		err := tx.UpdateColumns(post, "votes")
		if err != nil {
			return fmt.Errorf("upvoting failed: %v", err)
		}
		// TODO: What we need is an audit log.
		vote := &models.Vote{PostID: postId, UserID: u.ID}
		verrs, err := tx.ValidateAndCreate(vote)
		if err != nil {
			return fmt.Errorf("upvoting failed: %v", err)
		}
		if verrs.HasAny() {
			// TODO: This is an internal error.
			c.Set("errors", verrs)
			return c.Render(http.StatusUnprocessableEntity, r.JavaScript("error.js"))
		}
	} else {
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("error.js"))
	}

	c.Set("post", post)
	c.Set("board", &post.Board)
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

	post := &models.Post{}
	if err := tx.Eager().Find(post, postId); err != nil {
		return err
	}

	postVotes, downvoted, err := post.Board.Downvote(models.Redis, u, post)
	if err != nil {
		return err
	}
	if downvoted {
		post.Votes = postVotes
		err := tx.UpdateColumns(post, "votes")
		if err != nil {
			return fmt.Errorf("upvoting failed: %v", err)
		}
		// TODO: What we need is an audit log.
		// Indeed, it may be intereesting how often ppl change their votes.
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
	} else {
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("votes/fail.js"))
	}
	c.Set("post", post)
	c.Set("board", &post.Board)
	return c.Render(http.StatusOK, r.JavaScript("votes/create.js"))
}
