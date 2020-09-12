package actions

import (
	"github.com/gobuffalo/buffalo"
)

type UnauthenticatedController struct {
	Controller
}

func WithUnauthenticatedController(action func(c UnauthenticatedController) error) func(ctx buffalo.Context) error {
	return func(ctx buffalo.Context) error {
		c := UnauthenticatedController{}
		if err := c.SetUp(ctx); err != nil {
			return err
		}
		return action(c)
	}
}

func (c *UnauthenticatedController) SetUp(ctx buffalo.Context) error {
	if err := c.Controller.SetUp(ctx); err != nil {
		return err
	}

	return nil
}
