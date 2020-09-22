package actions

import (
	"fmt"
	"math/rand"
	"net/url"
	"os"
	"rally/models"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

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
			"userAvatarURL": func(size string, help plush.HelperContext) string {
				u, err := CurrentUser(help)
				if err != nil {
					return ""
				}
				return avatarURL(*u, size, true)
			},
			"isOriginalPosterComment": func(comment, post interface{}) bool {
				c := toCommentPtr(comment)
				p := toPostPtr(post)
				if c.Anonymous || p.Anonymous {
					return false
				}
				return c.AuthorID == p.AuthorID
			},
			"canManageComment": func(comment interface{}, help plush.HelperContext) bool {
				return canManageComment(comment, help.Context)
			},
			"canManagePost": func(post interface{}, help plush.HelperContext) bool {
				return canManagePost(post, help.Context) // Crashes otherwise.
			},
			"canManageBoard": func(board interface{}, help plush.HelperContext) bool {
				return canManageBoard(board, help.Context) // Crashes otherwise.
			},
			"isBoardStarred": func(board interface{}, help plush.HelperContext) bool {
				b := toBoardPtr(board)
				u, err := CurrentUser(help)
				if err != nil {
					log.Errorf("Unable to retrieve current user to determine if board starred: %v", err)
					return false
				}
				return u.IsBoardStarred(b)
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
		},
	})
}

func avatarURL(user models.User, size string, round bool) string {
	//
	// TODO: Resizing doesn't really work for Google default profile pictures
	// so we return an unresized url and <img> has width and height.
	if user.AvatarURL.Valid {
		cloudinaryID, ok := os.LookupEnv("CLOUDINARY_ACCOUNT")
		if !ok {
			log.Warn("CLOUDINARY_ACCOUNT environment variable not set, avatars will be ugly, resized by <img>")
			return user.AvatarURL.String
		}
		w := sizeToPixels(size)
		h := w
		cloudinaryURL := "https://res.cloudinary.com/%v/image/fetch/c_thumb,f_auto,fl_alpha,h_%v,r_30,w_%v/%v"
		return fmt.Sprintf(cloudinaryURL, cloudinaryID, h, w, user.AvatarURL.String)
	}

	return generateAvatarURL(user.Email, size, round)
}

func anonymousAvatarURL(size string, round bool) string {
	return generateAvatarURL(anonymousAvatarSeed, size, round)
}
func sizeToPixels(size string) int {
	if size == "large" {
		return 64
	} else if size == "small" {
		return 32
	} else {
		return 24
	}
}

func generateAvatarURL(seed, size string, round bool) string {
	px := sizeToPixels(size)
	var opts string
	if round {
		opts = "&b=%%23ffffff&m=8&r=50"
	}
	return fmt.Sprintf("https://avatars.dicebear.com/api/bottts/%v.com.svg?colorful=1&w=%v&h=%v&deterministic=1%v", url.QueryEscape(seed), px, px, opts)
}
