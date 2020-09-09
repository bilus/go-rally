package actions

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rally/models"
	"rally/services"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
	"github.com/pkg/errors"
)

//UsersNew renders the users form
func UsersNew(c buffalo.Context) error {
	u := models.User{}
	c.Set("user", u)
	return c.Render(200, r.HTML("users/new.plush.html"))
}

// UsersCreate registers a new user with the application.
func UsersCreate(c buffalo.Context) error {
	u := &models.User{}
	if err := c.Bind(u); err != nil {
		return errors.WithStack(err)
	}

	tx := c.Value("tx").(*pop.Connection)
	verrs, err := u.Create(tx)
	if err != nil {
		return errors.WithStack(err)
	}

	if verrs.HasAny() {
		c.Set("user", u)
		c.Set("errors", verrs)
		return c.Render(200, r.HTML("users/new.plush.html"))
	}

	return Login(u, c)
}

func UserDashboard(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	u, err := CurrentUser(c)
	if err != nil {
		return err
	}

	starredBoards, err := services.ListStarredBoards(tx, u, nil)
	if err != nil {
		return err
	}

	numBoards, err := tx.Count(&models.Board{})
	if err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("starredBoards", starredBoards)
		if numBoards == 0 {
			c.Set("showWelcome", true)
		}
		return c.Render(http.StatusOK, r.HTML("/users/dashboard.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(starredBoards))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(starredBoards))
	}).Respond(c)
}

func Login(u *models.User, c buffalo.Context) error {
	c.Session().Set("current_user_id", u.ID)

	redirectURL := "/"
	if redir, ok := c.Session().Get("redirectURL").(string); ok && redir != "" {
		redirectURL = redir
	}

	return c.Redirect(302, redirectURL)
}

// SetCurrentUser attempts to find a user based on the current_user_id
// in the session. If one is found it is set on the context.
func SetCurrentUser(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		RefreshCurrentUser(c)
		return next(c)
	}
}

// Authorize require a user be logged in before accessing a route
func Authorize(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		if uid := c.Session().Get("current_user_id"); uid == nil {
			c.Session().Set("redirectURL", c.Request().URL.String())

			err := c.Session().Save()
			if err != nil {
				return errors.WithStack(err)
			}

			c.Flash().Add("danger", "You must be authorized to see that page")
			return c.Redirect(302, "/auth/new")
		}
		return next(c)
	}
}

func RefreshCurrentUser(c buffalo.Context) {
	if uid := c.Session().Get("current_user_id"); uid != nil {
		u := &models.User{}
		tx := c.Value("tx").(*pop.Connection)
		err := tx.Find(u, uid)
		if err != nil {
			c.Session().Delete("current_user_id") // User deleted?
			c.Set("current_user", nil)
		}
		c.Set("current_user", u)
		boards, err := services.QuickAccessBoards(tx, u)
		if err != nil {
			log.Printf("Error loading quick access boards: %v", err)
			boards = []models.Board{}
		}
		c.Set("quickAccessBoards", boards)
	}
}

func CurrentUser(c context.Context) (*models.User, error) {
	u, ok := c.Value("current_user").(*models.User)
	if !ok {
		return nil, errors.New("current user in session has incorrect type")
	}
	return u, nil
}
