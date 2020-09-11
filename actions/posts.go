package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
	"github.com/gofrs/uuid"
)

// List gets all Posts.
// GET /posts
func (c AuthenticatedController) List() error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	posts := &models.Posts{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	drafts := c.Param("drafts") == "true"
	if drafts {
		currentUser, err := CurrentUser(c)
		if err != nil {
			return err
		}
		q = q.Where("draft AND author_id = ?", currentUser.ID)
	} else {
		q = q.Where("NOT draft")
	}

	// Can nest under /boards/:board_id
	boardId := c.Param("board_id")
	if boardId != "" {
		q = q.Where("board_id = ?", boardId)
	}

	order := c.Param("order")
	if order == "" {
		order = "top"
	}
	c.Set("orderClass", orderClass(order))
	if order == "newest" {
		q = q.Order("created_at DESC")
	} else {
		q = q.Order("votes DESC")
	}

	// Retrieve all Posts from the DB
	if err := q.Eager().All(posts); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("posts", posts)
		c.Set("drafts", drafts)
		return c.Render(http.StatusOK, r.HTML("posts/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(posts))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(posts))
	}).Respond(c)
}

func orderClass(activeOrder string) func(order string) string {
	return func(order string) string {
		if order == activeOrder {
			return "active"
		}
		return ""
	}
}

// Show gets the data for one Post. This function is mapped to
// the path GET /posts/{post_id}
func (c AuthenticatedController) Show() error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Post
	post := &models.Post{}

	// To find the Post the parameter post_id is used.
	if err := tx.Eager().Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("post", post)
		return c.Render(http.StatusOK, r.HTML("/posts/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(post))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(post))
	}).Respond(c)
}

// Create adds a Post to the DB. This function is mapped to the
// path POST /posts
func (c AuthenticatedController) Create() error {
	boardID, err := uuid.FromString(c.Param("board_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Allocate an empty Post
	post := &models.Post{Draft: true, BoardID: boardID}

	currentUser, err := CurrentUser(c)
	if err != nil {
		return err
	}
	post.AuthorID = currentUser.ID

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(post)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return fmt.Errorf("validation failed when creating an empty new draft: %q", verrs.String())
	}

	c.Set("post", post)

	return c.Render(http.StatusOK, r.JavaScript("/posts/create.plush.js"))
}

// Edit renders a edit form for a Post. This function is
// mapped to the path GET /posts/{post_id}/edit
func (c AuthenticatedController) Edit() error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Post
	post := &models.Post{}

	if err := tx.Eager().Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := authorizePostManagement(post, c); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	c.Set("post", post)
	return c.Render(http.StatusOK, r.HTML("/posts/edit.plush.html"))
}

// Update changes a Post in the DB. This function is mapped to
// the path PUT /posts/{post_id}
func (c AuthenticatedController) Update() error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Post
	post := &models.Post{}

	if err := tx.Eager().Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := authorizePostManagement(post, c); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	// Bind Post to the html form elements
	if err := c.Bind(post); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(post)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("post", post)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/posts/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "post.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/posts/%v", post.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(post))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(post))
	}).Respond(c)
}

// Destroy deletes a Post from the DB. This function is mapped
// to the path DELETE /posts/{post_id}
func (c AuthenticatedController) Destroy() error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Post
	post := &models.Post{}

	// To find the Post the parameter post_id is used.
	if err := tx.Eager().Find(post, c.Param("post_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := authorizePostManagement(post, c); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	if err := tx.Destroy(post); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "post.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/boards/%v", post.BoardID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(post))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(post))
	}).Respond(c)
}
