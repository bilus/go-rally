package actions

import (
	"os"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/logger"
	csrf "github.com/gobuffalo/mw-csrf"
	forcessl "github.com/gobuffalo/mw-forcessl"
	paramlogger "github.com/gobuffalo/mw-paramlogger"
	"github.com/markbates/goth/gothic"
	"github.com/unrolled/secure"

	"rally/models"

	"github.com/gobuffalo/buffalo-pop/v2/pop/popmw"
	i18n "github.com/gobuffalo/mw-i18n"
	"github.com/gobuffalo/packr/v2"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App
var T *i18n.Translator

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
//
// Routing, middleware, groups, etc... are declared TOP -> DOWN.
// This means if you add a middleware to `app` *after* declaring a
// group, that group will NOT have that new middleware. The same
// is true of resource declarations as well.
//
// It also means that routes are checked in the order they are declared.
// `ServeFiles` is a CATCH-ALL route, so it should always be
// placed last in the route declarations, as it will prevent routes
// declared after it to never be called.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:         ENV,
			SessionName: "_rally_session",
			Logger:      logger.New(logger.DebugLevel),
			Host:        getEnv("APP_HOST", ""),
		})

		// Automatically redirect to SSL
		app.Use(forceSSL())

		// Log request parameters (filters apply).
		app.Use(paramlogger.ParameterLogger)

		// Protect against CSRF attacks. https://www.owasp.org/index.php/Cross-Site_Request_Forgery_(CSRF)
		// Remove to disable this.
		csrf := csrf.New
		app.Use(csrf)

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.Connection)
		// Remove to disable this.
		app.Use(popmw.Transaction(models.DB))

		// Setup and use translations:
		app.Use(translations())

		//AuthMiddlewares
		app.Use(SetCurrentUser)
		app.Use(Authorize)

		postsResource := PostsResource{}
		app.GET("/", BoardsResource{}.List)

		app.POST("/posts", postsResource.Create)
		app.GET("/posts/{post_id}/edit", postsResource.Edit)
		app.GET("/posts/{post_id}", postsResource.Show)
		app.PUT("/posts/{post_id}", postsResource.Update)
		app.DELETE("/posts/{post_id}", postsResource.Destroy)

		app.Resource("/posts", PostsResource{})
		app.POST("/posts/{post_id}/votes", VotesCreate)
		app.DELETE("/posts/{post_id}/votes", VotesDestroy)
		app.POST("/posts/{post_id}/images", ImagesCreate)
		app.GET("/posts/{post_id}/images/{image_id}", ImagesShow)
		app.Middleware.Skip(csrf, ImagesCreate) // TODO: Handle csrf token sent by the editor.

		commentsResource := CommentsResource{}
		app.GET("/posts/{post_id}/comments", commentsResource.List)
		app.POST("/posts/{post_id}/comments", commentsResource.Create)
		app.DELETE("/posts/{post_id}/comments/{comment_id}", commentsResource.Destroy)

		//Routes for Auth
		auth := app.Group("/auth")
		auth.GET("/new", AuthNew)
		auth.POST("/", AuthCreate)
		auth.DELETE("/", AuthDestroy)
		authProviderNew := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)
		auth.GET("/{provider}", authProviderNew)
		auth.GET("/{provider}/callback", AuthCallback)
		auth.Middleware.Skip(Authorize, AuthNew, AuthCreate, authProviderNew, AuthCallback)

		if isSignupEnabled() {
			//Routes for User registration
			users := app.Group("/users")

			users.GET("/new", UsersNew)
			users.POST("/", UsersCreate)
			users.Middleware.Remove(Authorize)
		}

		app.Resource("/boards", BoardsResource{})
		app.ServeFiles("/", assetsBox) // serve files from the public directory
	}

	return app
}

func getEnv(name, def string) string {
	v, ok := os.LookupEnv(name)
	if ok {
		return v
	}
	return def
}

// translations will load locale files, set up the translator `actions.T`,
// and will return a middleware to use to load the correct locale for each
// request.
// for more information: https://gobuffalo.io/en/docs/localization
func translations() buffalo.MiddlewareFunc {
	var err error
	if T, err = i18n.New(packr.New("app:locales", "../locales"), "en-US"); err != nil {
		app.Stop(err)
	}
	return T.Middleware()
}

// forceSSL will return a middleware that will redirect an incoming request
// if it is not HTTPS. "http://example.com" => "https://example.com".
// This middleware does **not** enable SSL. for your application. To do that
// we recommend using a proxy: https://gobuffalo.io/en/docs/proxy
// for more information: https://github.com/unrolled/secure/
func forceSSL() buffalo.MiddlewareFunc {
	return forcessl.Middleware(secure.Options{
		SSLRedirect:     ENV == "production",
		SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
	})
}
