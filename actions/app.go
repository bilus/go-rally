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

		app.GET("/", WithAuthenticatedController(AuthenticatedController.Home))
		app.GET("/changelog", WithAuthenticatedController(AuthenticatedController.Changelog))
		app.GET("/dashboard", WithAuthenticatedController(AuthenticatedController.UserDashboard))

		// IDEA: Pass struct specifying route params, annotate required with required
		// Or make the struct the param and use reflection!
		app.GET("/boards/{board_id}/posts/{post_id}/edit", WithPostsController(PostsController.Edit))
		app.GET("/boards/{board_id}/posts/{post_id}", WithPostsController(PostsController.Show))
		app.PUT("/boards/{board_id}/posts/{post_id}", WithPostsController(PostsController.Update))
		app.DELETE("/boards/{board_id}/posts/{post_id}", WithPostsController(PostsController.Destroy))

		app.POST("/boards/{board_id}/posts/{post_id}/votes", WithPostsController(PostsController.VotesCreate))
		app.DELETE("/boards/{board_id}/posts/{post_id}/votes", WithPostsController(PostsController.VotesDestroy))

		// IMPORTANT: Buffalo Skip is stupid because it uses function name as the key.
		// But because a function returned from multiple invocations of WithPostsController has the same name every time,
		// Skipping it for one handler, skips it for all handlers to the given controller type.
		imagesCreate := func(ctx buffalo.Context) error { return WithPostsController(PostsController.ImagesCreate)(ctx) }
		app.POST("/posts/{post_id}/images", imagesCreate)
		app.GET("/posts/{post_id}/images/{image_id}", WithPostsController(PostsController.ImagesShow))

		app.Middleware.Skip(csrf, imagesCreate) // TODO: Handle csrf token sent by the editor.

		app.GET("/posts/{post_id}/comments", WithCommentsController(CommentsController.List))
		app.POST("/posts/{post_id}/comments", WithCommentsController(CommentsController.Create))
		app.DELETE("/posts/{post_id}/comments/{comment_id}", WithCommentsController(CommentsController.Destroy))

		//Routes for Auth
		auth := app.Group("/auth")
		authNew := WithUnauthenticatedController(UnauthenticatedController.AuthNew)
		auth.GET("/new", authNew)
		authCreate := WithUnauthenticatedController(UnauthenticatedController.AuthCreate)
		auth.POST("/", authCreate)
		authProviderNew := buffalo.WrapHandlerFunc(gothic.BeginAuthHandler)
		auth.GET("/{provider}", authProviderNew)
		authCallback := WithUnauthenticatedController(UnauthenticatedController.AuthCallback)
		auth.GET("/{provider}/callback", authCallback)
		// IMPORTANT: Skipping for any UnauthenticatedController action skips it for all the others.
		// See the note above.
		auth.Middleware.Skip(Authorize, authNew, authCreate, authProviderNew, authCallback)

		authDestroy := WithAuthenticatedController(AuthenticatedController.AuthDestroy)
		auth.DELETE("/", authDestroy)

		if isSignupEnabled() {
			//Routes for User registration
			users := app.Group("/users")

			users.GET("/new", WithUnauthenticatedController(UnauthenticatedController.UsersNew))
			users.POST("/", WithUnauthenticatedController(UnauthenticatedController.UsersCreate))
			// NOTE: The routes are unautherized anyway because they use UnauthenticatedController which
			// is referenced around the auth routes. See the notes above.
			users.Middleware.Remove(Authorize)
		}

		app.GET("/boards", WithBoardsController(BoardsController.List))
		app.GET("/boards/new", WithBoardsController(BoardsController.New))
		app.GET("/boards/{board_id}", WithBoardsController(BoardsController.Show))
		app.POST("/boards", WithBoardsController(BoardsController.Create))
		app.GET("/boards/{board_id}/edit", WithBoardsController(BoardsController.Edit))
		app.PUT("/boards/{board_id}", WithBoardsController(BoardsController.Update))
		app.DELETE("/boards/{board_id}", WithBoardsController(BoardsController.Destroy))
		app.POST("/boards/{board_id}/posts", WithPostsController(PostsController.Create))
		app.POST("/boards/{board_id}/refill", WithBoardsController(BoardsController.RefillCreate))
		app.POST("/boards/{board_id}/star", WithBoardsController(BoardsController.StarCreate))
		app.DELETE("/boards/{board_id}/star", WithBoardsController(BoardsController.StarDestroy))

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
