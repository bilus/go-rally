package actions

import (
	"rally/models"

	"github.com/gobuffalo/buffalo"
)

type AuthenticatedController struct {
	Controller

	CurrentUser models.User
}

func WithAuthenticatedController(action func(c AuthenticatedController) error) func(c buffalo.Context) error {
	return func(ctx buffalo.Context) error {
		c := AuthenticatedController{}
		if err := c.SetUp(ctx); err != nil {
			return err
		}
		return action(c)
	}
}

func (c *AuthenticatedController) SetUp(ctx buffalo.Context) error {
	if err := c.Controller.SetUp(ctx); err != nil {
		return err
	}
	user, err := CurrentUser(ctx)
	if err != nil {
		return err
	}
	c.CurrentUser = *user
	return nil
}
