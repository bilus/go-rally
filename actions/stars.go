package actions

import (
	"net/http"
	"rally/services"
)

func (c BoardsController) StarCreate() error {
	if err := c.RequireBoard(); err != nil {
		return err
	}

	if err := services.StarBoard(c.Tx, &c.CurrentUser, c.Board, true); err != nil {
		return c.Error(http.StatusUnprocessableEntity, err)
	}

	c.Set("board", c.Board)
	return c.Render(http.StatusOK, r.JavaScript("/stars/create.plush.js"))
}

func (c BoardsController) StarDestroy() error {
	if err := c.RequireBoard(); err != nil {
		return err
	}

	if err := services.StarBoard(c.Tx, &c.CurrentUser, c.Board, false); err != nil {
		return c.Error(http.StatusUnprocessableEntity, err)
	}

	c.Set("board", c.Board)
	return c.Render(http.StatusOK, r.JavaScript("/stars/destroy.plush.js"))
}
