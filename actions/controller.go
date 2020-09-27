package actions

import (
	"fmt"
	"math/rand"
	"net/http"
	"rally/models"
	"rally/services"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/flect"
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

	c.RegisterHelpers(c.baseHelpers())

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

func (c *Controller) RegisterHelpers(h render.Helpers) {
	for k, v := range h {
		c.Set(k, v)
	}
}

func (c *Controller) baseHelpers() render.Helpers {
	return render.Helpers{
		"timeAgo": func(t time.Time) string {
			return humanize.Time(t)
		},
		"safeTitle": func(t string) string {
			if t == "" {
				return untitled
			} else {
				return t
			}
		},
		"postAuthor": func(post interface{}) string {
			p := toPostPtr(post)

			if p.Anonymous {
				return anonymousName
			}

			return p.Author.Email
		},
		"postAvatarURL": func(post interface{}, size string) string {
			if size == "" {
				size = "large"
			}
			p := toPostPtr(post)
			if p.Anonymous {
				return anonymousAvatarURL(size, false)
			}
			return avatarURL(p.Author, size, false)
		},
		"commentAuthor": func(comment interface{}) string {
			c := toCommentPtr(comment)
			if c.Anonymous {
				return anonymousName
			}
			return c.Author.Email
		},
		"commentAvatarURL": func(comment interface{}) string {
			c := toCommentPtr(comment)
			if c.Anonymous {
				return anonymousAvatarURL("small", false)
			}
			return avatarURL(c.Author, "small", false)
		},
		"avatarURL": avatarURL,
		"maybePluralize": func(noun string, count int) string {
			if count == 1 {
				return noun
			}
			return flect.Pluralize(noun)
		},
		"isOriginalPosterComment": func(comment, post interface{}) bool {
			c := toCommentPtr(comment)
			p := toPostPtr(post)
			if c.Anonymous || p.Anonymous {
				return false
			}
			return c.AuthorID == p.AuthorID
		},
		"greeting": func() string {
			greetings := []string{
				"Salut! 👋", "Hola! 👋", "Privet! 👋", "Nǐ hǎo 👋", "Ciao! 👋", "Yā, Yō 👋", "Hallo! 👋", "Hello! 👋", "Cześć! 👋", "Helo! 👋", "Hey! 👋", "Hei! 👋",
			}
			return greetings[rand.Intn(len(greetings))]
		},
		"formatUserList": func(users []models.User) string {
			emails := make([]string, len(users))
			for i, user := range users {
				emails[i] = user.Email
			}
			return strings.Join(emails, "<br />")
		},
	}
}
