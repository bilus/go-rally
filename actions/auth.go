package actions

import (
	"fmt"
	"log"
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
)

func init() {
	gothic.Store = App().SessionStore

	goth.UseProviders(
		google.New(mustGetEnv("GOOGLE_KEY"), mustGetEnv("GOOGLE_SECRET"), fmt.Sprintf("%s%s", App().Host, "/auth/google/callback")),
	)
}

func mustGetEnv(name string) string {
	v, ok := os.LookupEnv(name)
	if !ok {
		log.Fatalf("Required environment variable missing: %q", name)
	}
	return v
}

func AuthCallback(c buffalo.Context) error {
	user, err := gothic.CompleteUserAuth(c.Response(), c.Request())
	if err != nil {
		return c.Error(401, err)
	}
	// Do something with the user, maybe register them/sign them in
	return c.Render(200, r.JSON(user))
}
