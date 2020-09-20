package actions

import (
	"fmt"
	"net/http"
)

func (c PostsController) CreateReaction() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	emoji := c.Param("emoji")
	if emoji == "" {
		err := fmt.Errorf("missing emoji param")
		c.Set("error", err)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("/fail.js"))
	}

	if err := c.ReactionsService.AddReactionToPost(&c.CurrentUser, c.Post, emoji); err != nil {
		c.Set("error", err)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("/fail.js"))
	}

	return c.Render(http.StatusCreated, r.JavaScript("reactions/create.js"))
}

func (c PostsController) DestroyReaction() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	emoji := c.Param("emoji")
	if emoji == "" {
		err := fmt.Errorf("missing emoji param")
		c.Set("error", err)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("/fail.js"))
	}

	if err := c.ReactionsService.RemoveReactionToPost(&c.CurrentUser, c.Post, emoji); err != nil {
		c.Set("error", err)
		return c.Render(http.StatusUnprocessableEntity, r.JavaScript("/fail.js"))
	}

	return c.Render(http.StatusCreated, r.JavaScript("reactions/destroy.js"))
}
