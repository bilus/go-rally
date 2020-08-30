package actions

import (
	"context"
	"fmt"
	"net/url"
	"rally/models"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/flect"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/plush/v4"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

const anonymousName = "Anonymous"
const anonymousAvatarSeed = "anonymous@example.com"
const untitled = "Untitled"

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.plush.html",

		// Box containing all of the templates:
		TemplatesBox: packr.New("app:templates", "../templates"),
		AssetsBox:    assetsBox,

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
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
			"postAuthor": func(post *models.Post) string {
				if post.Anonymous {
					return anonymousName
				}

				return post.Author.Email
			},
			"postAvatarURL": func(post *models.Post) string {
				if post.Anonymous {
					return avatarURL(anonymousAvatarSeed, "large")
				}
				return avatarURL(post.Author.Email, "large")
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
					return avatarURL(anonymousAvatarSeed, "small")
				}
				return avatarURL(c.Author.Email, "small")
			},
			"avatarURL": avatarURL,
			"maybePluralize": func(noun string, count int) string {
				if count == 1 {
					return noun
				}
				return flect.Pluralize(noun)
			},
			"canManageComment": func(comment interface{}, help plush.HelperContext) bool {
				return canManageComment(comment, help.Context)
			},
			"canManagePost": func(post interface{}, help plush.HelperContext) bool {
				return canManagePost(post, help.Context) // Crashes otherwise.
			},
		},
	})
}

func avatarURL(seed, size string) string {
	var px int
	if size == "large" {
		px = 64
	} else {
		px = 32
	}
	return fmt.Sprintf("https://avatars.dicebear.com/api/bottts/%v.com.svg?colorful=1&w=%v&h=%v&deterministic=1", url.QueryEscape(seed), px, px)
}

func toCommentPtr(comment interface{}) *models.Comment {
	ptr, ok := comment.(*models.Comment)
	if ok {
		return ptr
	}
	val, ok := comment.(models.Comment)
	if ok {
		return &val
	}
	panic("Expecting models.Comment or *models.Comment")
}

func toPostPtr(post interface{}) *models.Post {
	ptr, ok := post.(*models.Post)
	if ok {
		return ptr
	}
	val, ok := post.(models.Post)
	if ok {
		return &val
	}
	panic("Expecting models.Post or *models.Post")
}

func canManageComment(comment interface{}, help context.Context) bool {
	c := toCommentPtr(comment)
	u, err := CurrentUser(help)
	if err != nil {
		return false
	}
	return c.AuthorID == u.ID
}

func canManagePost(post interface{}, help context.Context) bool {
	p := toPostPtr(post)
	u, err := CurrentUser(help)
	if err != nil {
		return false
	}
	return p.AuthorID == u.ID
}
