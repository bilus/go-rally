package grifts

import (
	"rally/models"

	"github.com/gobuffalo/pop/v5"
	"github.com/markbates/grift/grift"
)

var _ = grift.Namespace("db", func() {

	grift.Desc("seed", "Seeds a database")
	grift.Add("seed", func(c *grift.Context) error {
		return models.DB.Transaction(func(tx *pop.Connection) error {
			password := "12345678"

			u := &models.User{
				Email:                "someone@example.com",
				Password:             password,
				PasswordConfirmation: password,
			}

			_, err := u.Create(tx)
			if err != nil {
				return err
			}

			return nil
		})
		return nil
	})

})
