package actions

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strings"

	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

func init() {
	gothic.Store = App().SessionStore

	google := google.New(os.Getenv("GOOGLE_KEY"), os.Getenv("GOOGLE_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/google/callback"))
	googleHostedDomain := os.Getenv("GOOGLE_HOSTED_DOMAIN")
	if len(googleHostedDomain) > 0 {
		google.SetHostedDomain(googleHostedDomain)
	}
	goth.UseProviders(google)
}

func AuthCallback(c buffalo.Context) error {
	profile, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	tx := c.Value("tx").(*pop.Connection)
	q := tx.Where("email = ?", profile.Email)
	user := models.User{}
	err = q.First(&user)
	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			user.Email = profile.Email
			user.GoogleUserID = profile.UserID // TODO: Extend to other providers.
			verrs, err := tx.ValidateAndCreate(&user)
			if err != nil {
				return err
			}
			if verrs.Count() > 0 {
				c.Set("errors", verrs)
				c.Set("user", &user)

				return c.Render(http.StatusUnauthorized, r.HTML("auth/new.plush.html"))
			}

		} else {
			return err
		}
	}

	return Login(&user, c)
}

// AuthNew loads the signin page
func AuthNew(c buffalo.Context) error {
	c.Set("user", models.User{})
	return c.Render(200, r.HTML("auth/new.plush.html"))
}

// AuthCreate attempts to log the user in with an existing account.
func AuthCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)

	// find a user with the email
	err := tx.Where("email = ?", strings.ToLower(strings.TrimSpace(u.Email))).First(u)

	// helper function to handle bad attempts
	bad := func() error {
		verrs := validate.NewErrors()
		verrs.Add("email", "invalid email/password")

		c.Set("errors", verrs)
		c.Set("user", u)

		return c.Render(http.StatusUnauthorized, r.HTML("auth/new.plush.html"))
	}

	if err != nil {
		if errors.Cause(err) == sql.ErrNoRows {
			// couldn't find an user with the supplied email address.
			return bad()
		}
		return errors.WithStack(err)
	}

	// confirm that the given password matches the hashed password from the db
	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		return bad()
	}

	return Login(u, c)
}

// AuthDestroy clears the session and logs a user out
func AuthDestroy(c buffalo.Context) error {
	c.Session().Clear()
	c.Flash().Add("success", "You have been logged out!")
	return c.Redirect(302, "newAuthPath()")
}
