package actions

import (
	"fmt"
	"net/http"
	"rally/models"
	"rally/services"
	"rally/stores"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

type PostsController struct {
	BoardsController

	Post *models.Post

	ReactionsService services.ReactionsService
}

func WithPostsController(action func(c PostsController) error) func(ctx buffalo.Context) error {
	return func(ctx buffalo.Context) error {
		c := PostsController{}
		if err := c.SetUp(ctx); err != nil {
			return err
		}
		return action(c)
	}
}

func (c *PostsController) SetUp(ctx buffalo.Context) error {
	if err := c.BoardsController.SetUp(ctx); err != nil {
		return err
	}

	postIDParam := ctx.Param("post_id")
	if postIDParam != "" {
		postID, err := uuid.FromString(postIDParam)
		if err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}

		c.Post = &models.Post{}
		if err := c.Tx.Eager().Find(c.Post, postID); err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}
		c.ReactionsService = services.NewReactionsService(stores.NewReactionStore(models.Redis), c.Tx)
	}

	return nil
}

// TODO: Allow specifying interceptors in With..Controller and use them.
func (c *PostsController) RequirePost() error {
	if c.Post == nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("missing post ID"))
	}
	return nil
}

func (c *PostsController) RequirePostWithWriteAccess() error {
	if err := c.RequirePost(); err != nil {
		return err
	}
	if err := authorizePostManagement(c.Post, c); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	return nil
}
