package actions

import (
	"fmt"
	"net/http"
	"rally/models"
	"rally/services"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

// TODO: Lots of duplication, different 2 lines.

func StarCreate(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	u, err := CurrentUser(c)
	if err != nil {
		return err
	}

	boardID, err := uuid.FromString(c.Param("board_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	board := &models.Board{}
	if err := tx.Find(board, boardID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := services.StarBoard(tx, u, board, true); err != nil {
		return c.Error(http.StatusUnprocessableEntity, err)
	}

	c.Set("board", board)

	return c.Render(http.StatusOK, r.JavaScript("/stars/create.plush.js"))
}

func StarDestroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	u, err := CurrentUser(c)
	if err != nil {
		return err
	}

	boardID, err := uuid.FromString(c.Param("board_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	board := &models.Board{}
	if err := tx.Find(board, boardID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := services.StarBoard(tx, u, board, false); err != nil {
		return c.Error(http.StatusUnprocessableEntity, err)
	}

	c.Set("board", board)

	return c.Render(http.StatusOK, r.JavaScript("/stars/destroy.plush.js"))
}
