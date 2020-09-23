package actions

import (
	"fmt"
	"net/http"
	"rally/services"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
)

type Controller struct {
	buffalo.Context
	Tx               *pop.Connection
	PaginationParams pop.PaginationParams
}

func WithController(action func(c Controller) error) func(ctx buffalo.Context) error {
	return func(c buffalo.Context) error {
		ct := Controller{}
		if err := ct.SetUp(c); err != nil {
			return err
		}
		return ct.wrapServiceError(action(ct))
	}
}

func (c *Controller) SetUp(ctx buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := ctx.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}
	if tx == nil {
		return fmt.Errorf("no transaction found")
	}
	c.Tx = tx
	c.Context = ctx
	c.PaginationParams = ctx.Params()
	return nil
}

func (c *Controller) wrapServiceError(err error) error {
	if err == services.ErrNotFound {
		return c.Error(http.StatusNotFound, err)
	}
	if err == services.ErrUnauthorized {
		return c.Error(http.StatusUnauthorized, err)
	}
	return err
}
