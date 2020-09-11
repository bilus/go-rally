package actions

import (
	"rally/models"

	"github.com/gobuffalo/buffalo"
)

type AuthenticatedController struct {
	Controller

	CurrentUser models.User
}

func (ct *AuthenticatedController) SetUp(c buffalo.Context) error {
	if err := ct.Controller.SetUp(c); err != nil {
		return err
	}
	user, err := CurrentUser(c)
	if err != nil {
		return err
	}
	ct.CurrentUser = *user
	return nil
}

func WithAuthenticatedController(action func(ct AuthenticatedController) error) func(c buffalo.Context) error {
	return func(c buffalo.Context) error {
		ct := AuthenticatedController{}
		if err := ct.SetUp(c); err != nil {
			return err
		}
		return action(ct)
	}
}
