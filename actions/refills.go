package actions

import (
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/x/responder"
)

func (c BoardsController) RefillCreate() error {
	if err := c.RequireBoardWithWriteAccess(); err != nil {
		return err
	}

	// TODO: Performance issue with large number of users.
	users := []models.User{}
	if err := c.Tx.RawQuery("SELECT id FROM users").Select("id").All(&users); err != nil {
		return err
	}

	// TODO: Show error.
	if err := c.VotingService.Refill(c.Board, users...); err != nil {
		return c.Error(http.StatusUnprocessableEntity, err)
	}
	return responder.Wants("javascript", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JavaScript("/refills/create.plush.js"))
	}).Respond(c)
}
