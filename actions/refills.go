package actions

import (
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/x/responder"
	"github.com/gofrs/uuid"
)

func (c BoardsController) RefillCreate() error {
	if err := c.RequireBoardWithWriteAccess(); err != nil {
		return err
	}

	// TODO: Performance issue with large number of users.
	users := []struct{ ID uuid.UUID }{}
	if err := c.Tx.RawQuery("SELECT id FROM users").All(&users); err != nil {
		return err
	}
	ids := make([]uuid.UUID, len(users))
	for i, u := range users {
		ids[i] = u.ID
	}

	// TODO: Show error.
	if err := c.Board.VotingStrategy.Refill(models.Redis, c.Board, ids...); err != nil {
		return c.Error(http.StatusUnprocessableEntity, err)
	}
	return responder.Wants("javascript", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JavaScript("/refills/create.plush.js"))
	}).Respond(c)
}
