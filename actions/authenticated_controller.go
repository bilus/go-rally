package actions

import (
	"rally/models"
	"rally/services"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

type AuthenticatedController struct {
	Controller

	CurrentUser      models.User
	DashboardService services.DashboardService
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
	c.DashboardService = services.NewDashboardService(c.Tx)

	c.Set("currentUser", *user)

	c.RegisterHelpers(render.Helpers{
		"userAvatarURL": func(size string) string {
			return avatarURL(c.CurrentUser, size, true)
		},
	})

	return nil
}
