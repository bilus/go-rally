package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
	"github.com/gofrs/uuid"
)

func RefillCreate(c buffalo.Context) error {
	boardID, err := uuid.FromString(c.Param("board_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Board
	board := &models.Board{}

	// To find the Board the parameter board_id is used.
	if err := tx.Find(board, boardID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := authorizeBoardManagement(board, c); err != nil {
		return err
	}

	// TODO: Performance issue with large number of users.
	users := []struct{ ID uuid.UUID }{}
	if err := tx.RawQuery("SELECT id FROM users").All(&users); err != nil {
		return err
	}
	ids := make([]uuid.UUID, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}

	// TODO: Show error.
	err = board.VotingStrategy.Refill(models.Redis, board, ids...)
	if err != nil {
		return c.Error(http.StatusUnprocessableEntity, err)
	}
	return responder.Wants("javascript", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JavaScript("/refills/create.plush.js"))
	}).Respond(c)
}
