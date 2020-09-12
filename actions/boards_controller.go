package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

type BoardsController struct {
	AuthenticatedController

	Board *models.Board
}

func WithBoardsController(action func(c BoardsController) error) func(c buffalo.Context) error {
	return func(ctx buffalo.Context) error {
		c := BoardsController{}
		if err := c.SetUp(ctx); err != nil {
			return err
		}
		return action(c)
	}
}

func (c *BoardsController) SetUp(ctx buffalo.Context) error {
	if err := c.AuthenticatedController.SetUp(ctx); err != nil {
		return err
	}

	boardIDParam := ctx.Param("board_id")
	if boardIDParam != "" {
		boardID, err := uuid.FromString(boardIDParam)
		if err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}
		c.Board = &models.Board{}
		if err := c.Tx.Find(c.Board, boardID); err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}
	}
	return nil
}

func (c *BoardsController) RequireBoard() error {
	if c.Board == nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("missing board ID"))
	}
	return nil
}

func (c *BoardsController) RequireBoardWithWriteAccess() error {
	if err := c.RequireBoard(); err != nil {
		return err
	}
	if err := authorizeBoardManagement(c.Board, c); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	return nil
}
