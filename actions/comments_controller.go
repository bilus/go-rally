package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

type CommentsController struct {
	PostsController

	Comment *models.Comment
}

func WithCommentsController(action func(c CommentsController) error) func(ctx buffalo.Context) error {
	return func(ctx buffalo.Context) error {
		c := CommentsController{}
		if err := c.SetUp(ctx); err != nil {
			return err
		}
		return action(c)
	}
}

func (c *CommentsController) SetUp(ctx buffalo.Context) error {
	if err := c.PostsController.SetUp(ctx); err != nil {
		return err
	}

	commentIDParam := ctx.Param("comment_id")
	if commentIDParam != "" {
		commentID, err := uuid.FromString(commentIDParam)
		if err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}
		c.Comment = &models.Comment{}
		if err := c.Tx.Eager().Find(c.Comment, commentID); err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}
	}
	return nil
}

// TODO: Allow specifying interceptors in With..Controller and use them?
func (c *CommentsController) RequireComment() error {
	if c.Comment == nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("missing comment ID"))
	}
	return nil
}

func (c *CommentsController) RequireCommentWithWriteAccess() error {
	if err := c.RequireComment(); err != nil {
		return err
	}
	if err := c.authorizeCommentManagement(c.Comment); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	return nil
}
