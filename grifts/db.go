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

			_, err := tx.ValidateAndSave(u)
			if err != nil {
				return err
			}

			p := &models.Post{
				Title:    "Lorem ipsum",
				Body:     "Lorem ipsum dolor sit amet, consectetuer adipiscing elit.  Donec hendrerit tempor tellus.  Donec pretium posuere tellus.  Proin quam nisl, tincidunt et, mattis eget, convallis nec, purus.  Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus.  Nulla posuere.  Donec vitae dolor.  Nullam tristique diam non turpis.  Cras placerat accumsan nulla.  Nullam rutrum.  Nam vestibulum accumsan nisl.",
				AuthorID: u.ID,
			}
			_, err = tx.ValidateAndSave(p)
			if err != nil {
				return err
			}
			return nil
		})
		return nil
	})

})
