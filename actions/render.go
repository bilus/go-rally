package actions

import (
	"rally/models"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/packr/v2"
)

var r *render.Engine
var assetsBox = packr.New("app:assets", "../public")

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
					return "Untitled"
				} else {
					return t
				}
			},
			"postAuthor": func(post *models.Post) string {
				if post.Anonymous {
					return "Anonymous"
				}

				return post.Author.Email
			},
		},
	})
}
