package actions

import (
	"fmt"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

type Controller struct {
	Tx *pop.Connection
}

func (ct *Controller) SetUp(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	ct.Tx = tx
	return nil
}

func WithController(action func(ct Controller, c buffalo.Context) error) func(c buffalo.Context) error {
	return func(c buffalo.Context) error {
		ct := Controller{}
		if err := ct.SetUp(c); err != nil {
			return err
		}
		return action(ct, c)
	}
}
