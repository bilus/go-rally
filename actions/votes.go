package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/pop/v5/slices"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
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
		logVotingAuditEvent(c, "upvote", post)
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
		logVotingAuditEvent(c, "downpvote", post)
	} else {
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("votes/fail.js"))
	}
	c.Set("post", post)
	c.Set("board", &post.Board)
	return c.Render(http.StatusOK, r.JavaScript("votes/destroy.js"))
}

func logVotingAuditEvent(c buffalo.Context, type_ string, post *models.Post) {
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		log.Errorf("unable to log event: %v", fmt.Errorf("no transaction found"))
	}

	userID := uuid.UUID{}
	u, err := CurrentUser(c)
	if err != nil {
		log.Errorf("unable to retrieve current user for auditing purposes: %v", err)
	} else {
		userID = u.ID
	}

	event := &models.AuditEvent{
		Type: type_,
		Payload: slices.Map{
			"current_user": userID,
			"board_id":     post.BoardID,
			"post_id":      post.ID,
			"author_id":    post.AuthorID,
			"post_votes":   post.Votes,
		},
	}
	err = tx.Create(event)
	if err != nil {
		log.Errorf("unable to log event: %v", err)
	}
}
