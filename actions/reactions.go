package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/pop/v5/slices"
	log "github.com/sirupsen/logrus"
)

func (c PostsController) CreateReaction() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	emoji := c.Param("emoji")
	if emoji == "" {
		err := fmt.Errorf("missing emoji param")
		c.Set("error", err)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("/fail.plush.js"))
	}

	if err := c.ReactionsService.AddReactionToPost(&c.CurrentUser, c.Post, emoji); err != nil {
		c.Set("error", err)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("/fail.plush.js"))
	}

	reactions, err := c.ReactionsService.ListReactionsToPost(&c.CurrentUser, c.Post)
	if err != nil {
		return err
	}
	c.Post.Reactions = reactions

	c.logReactionAuditEvent("add_reaction", emoji)

	c.Set("post", c.Post)
	return c.Render(http.StatusCreated, r.JavaScript("reactions/create.plush.js"))
}

func (c PostsController) DestroyReaction() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	emoji := c.Param("emoji")
	if emoji == "" {
		err := fmt.Errorf("missing emoji param")
		c.Set("error", err)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("/fail.plush.js"))
	}

	if err := c.ReactionsService.RemoveReactionToPost(&c.CurrentUser, c.Post, emoji); err != nil {
		c.Set("error", err)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("/fail.plush.js"))
	}

	reactions, err := c.ReactionsService.ListReactionsToPost(&c.CurrentUser, c.Post)
	if err != nil {
		return err
	}
	c.Post.Reactions = reactions

	c.logReactionAuditEvent("remove_reaction", emoji)

	c.Set("post", c.Post)
	return c.Render(http.StatusCreated, r.JavaScript("reactions/destroy.plush.js"))
}

func (c PostsController) logReactionAuditEvent(type_, emoji string) {
	event := &models.AuditEvent{
		Type: type_,
		Payload: slices.Map{
			"current_user": c.CurrentUser.ID,
			"board_id":     c.Post.BoardID,
			"post_id":      c.Post.ID,
			"author_id":    c.Post.AuthorID,
			"reaction":     emoji,
		},
	}
	if err := c.Tx.Create(event); err != nil {
		log.Errorf("unable to log event: %v", err)
	}
}
