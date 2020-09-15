package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/pop/v5/slices"
	log "github.com/sirupsen/logrus"
)

// VotesCreate upvotes a post
func (c PostsController) VotesCreate() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	postVotes, upvoted, err := c.VotingService.Upvote(&c.CurrentUser, c.Post)
	if err != nil {
		return err
	}

	if upvoted {
		c.Post.Votes = postVotes
		err := c.Tx.UpdateColumns(c.Post, "votes")
		if err != nil {
			return fmt.Errorf("upvoting failed: %v", err)
		}
		c.logVotingAuditEvent("upvote")
	} else {
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("error.js"))
	}

	c.Set("post", c.Post)
	c.Set("board", &c.Post.Board)
	return c.Render(http.StatusOK, r.JavaScript("votes/create.js"))
}

// VotesDestroy downvotes a post
func (c PostsController) VotesDestroy() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	postVotes, downvoted, err := c.VotingService.Downvote(&c.CurrentUser, c.Post)
	if err != nil {
		return err
	}
	if downvoted {
		c.Post.Votes = postVotes
		err := c.Tx.UpdateColumns(c.Post, "votes")
		if err != nil {
			return fmt.Errorf("upvoting failed: %v", err)
		}
		c.logVotingAuditEvent("downpvote")
	} else {
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("votes/fail.js"))
	}
	c.Set("post", c.Post)
	c.Set("board", &c.Post.Board)
	return c.Render(http.StatusOK, r.JavaScript("votes/destroy.js"))
}

func (c PostsController) logVotingAuditEvent(type_ string) {
	event := &models.AuditEvent{
		Type: type_,
		Payload: slices.Map{
			"current_user": c.CurrentUser.ID,
			"board_id":     c.Post.BoardID,
			"post_id":      c.Post.ID,
			"author_id":    c.Post.AuthorID,
			"post_votes":   c.Post.Votes,
		},
	}
	if err := c.Tx.Create(event); err != nil {
		log.Errorf("unable to log event: %v", err)
	}
}
