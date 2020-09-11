package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

type Controller struct {
	buffalo.Context
	Tx *pop.Connection
}

func WithController(action func(c Controller) error) func(ctx buffalo.Context) error {
	return func(c buffalo.Context) error {
		ct := Controller{}
		if err := ct.SetUp(c); err != nil {
			return err
		}
		return action(ct)
	}
}

func (c *Controller) SetUp(ctx buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := ctx.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	c.Tx = tx
	c.Context = ctx
	return nil
}
