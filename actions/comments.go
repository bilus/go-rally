package actions

import (
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
	"github.com/gofrs/uuid"
)

// List gets all Comments. This function is mapped to the path
// GET /comments
func (c CommentsController) List() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	comments := &models.Comments{}
	if err := listComments(c.Tx.Q(), c.Post.ID, comments); err != nil {
		return err
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// Paginate results. Params "page" and "per_page" control pagination.
		// Default values are "page=1" and "per_page=20".

		// Add the paginator to the context so it can be used in the template.
		ctx.Set("comments", comments)
		ctx.Set("post", c.Post)
		return ctx.Render(http.StatusOK, r.HTML("/comments/index.plush.html"))
	}).Wants("javascript", func(ctx buffalo.Context) error {
		ctx.Set("comments", comments)
		ctx.Set("comment", &models.Comment{PostID: c.Post.ID}) // Inline new comment form.
		ctx.Set("post", c.Post)
		return ctx.Render(http.StatusOK, r.JavaScript("/comments/index.plush.js"))
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.JSON(comments))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.XML(comments))
	}).Respond(c)
}

// Show gets the data for one Comment. This function is mapped to
// the path GET /comments/{comment_id}
func (c CommentsController) Show() error {
	if err := c.RequireComment(); err != nil {
		return err
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		ctx.Set("comment", c.Comment)
		return ctx.Render(http.StatusOK, r.HTML("/comments/show.plush.html"))
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.JSON(c.Comment))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.XML(c.Comment))
	}).Respond(c)
}

// New renders the form for creating a new Comment.
// This function is mapped to the path GET /comments/new
func (c CommentsController) New() error {
	c.Set("comment", &models.Comment{})

	return c.Render(http.StatusOK, r.HTML("/comments/new.plush.html"))
}

// Create adds a Comment to the DB. This function is mapped to the
// path POST /comments
func (c CommentsController) Create() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	// Bind comment to the html form elements
	comment := &models.Comment{}
	if err := c.Bind(comment); err != nil {
		return err
	}

	comment.PostID = c.Post.ID
	comment.AuthorID = c.CurrentUser.ID
	comment.Author = c.CurrentUser

	// Validate the data from the html form
	verrs, err := c.Tx.ValidateAndCreate(comment)
	if err != nil {
		return err
	}

	err = c.Tx.RawQuery("UPDATE posts SET comment_count = comment_count + 1 WHERE id = ?", comment.PostID).Exec()
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(ctx buffalo.Context) error {
			ctx.Set("errors", verrs)
			ctx.Set("comment", comment)
			return ctx.Render(http.StatusUnprocessableEntity, r.HTML("/comments/new.plush.html"))
		}).Wants("javascript", func(ctx buffalo.Context) error {
			// Make the errors available inside the html template
			ctx.Set("errors", verrs)
			ctx.Set("comment", comment)
			return ctx.Render(http.StatusUnprocessableEntity, r.JavaScript("/comments/failed.plush.js")) // TODO:
		}).Wants("json", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(cctx buffalo.Context) error {
			return cctx.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a success message
		ctx.Flash().Add("success", T.Translate(ctx, "comment.created.success"))

		// and redirect to the show page
		return ctx.Redirect(http.StatusSeeOther, "/comments/%v", comment.ID)
	}).Wants("javascript", func(ctx buffalo.Context) error {
		ctx.Set("comment", comment)
		ctx.Set("post", c.Post)

		comments := &models.Comments{}
		if err := listComments(pop.Q(c.Tx), c.Post.ID, comments); err != nil {
			return err
		}
		ctx.Set("comments", comments)

		return ctx.Render(http.StatusCreated, r.JavaScript("comments/create.plush.js")) // TODO:
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusCreated, r.JSON(comment))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusCreated, r.XML(comment))
	}).Respond(c)
}

// Edit renders a edit form for a Comment. This function is
// mapped to the path GET /comments/{comment_id}/edit
func (c CommentsController) Edit() error {
	if err := c.RequireCommentWithWriteAccess(); err != nil {
		return err
	}

	c.Set("comment", c.Comment)
	return c.Render(http.StatusOK, r.HTML("/comments/edit.plush.html"))
}

// Update changes a Comment in the DB. This function is mapped to
// the path PUT /comments/{comment_id}
func (c CommentsController) Update() error {
	if err := c.RequireCommentWithWriteAccess(); err != nil {
		return err
	}

	// Bind Comment to the html form elements
	if err := c.Bind(c.Comment); err != nil {
		return err
	}

	verrs, err := c.Tx.ValidateAndUpdate(c.Comment)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(ctx buffalo.Context) error {
			// Make the errors available inside the html template
			ctx.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			ctx.Set("comment", c.Comment)

			return ctx.Render(http.StatusUnprocessableEntity, r.HTML("/comments/edit.plush.html"))
		}).Wants("json", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a success message
		ctx.Flash().Add("success", T.Translate(ctx, "comment.updated.success"))

		// and redirect to the show page
		return ctx.Redirect(http.StatusSeeOther, "/comments/%v", c.Comment.ID)
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JSON(c.Comment))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.XML(c.Comment))
	}).Respond(c)
}

// Destroy deletes a Comment from the DB. This function is mapped
// to the path DELETE /comments/{comment_id}
func (c CommentsController) Destroy() error {
	if err := c.RequireCommentWithWriteAccess(); err != nil {
		return err
	}

	if err := c.Tx.Destroy(c.Comment); err != nil {
		return err
	}

	err := c.Tx.RawQuery("UPDATE posts SET comment_count = comment_count - 1 WHERE id = ?", c.Comment.PostID).Exec()
	if err != nil {
		return err
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a flash message
		ctx.Flash().Add("success", T.Translate(ctx, "comment.destroyed.success"))

		// Redirect to the index page
		return ctx.Redirect(http.StatusSeeOther, "/comments")
	}).Wants("javascript", func(ctx buffalo.Context) error {
		ctx.Set("post", c.Post)
		ctx.Set("comment", c.Comment)
		comments := &models.Comments{}
		if err := listComments(pop.Q(c.Tx), c.Post.ID, comments); err != nil {
			return err
		}
		ctx.Set("comments", comments)
		return ctx.Render(http.StatusOK, r.JavaScript("comments/destroy.plush.js"))
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JSON(c.Comment))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.XML(c.Comment))
	}).Respond(c)
}

func listComments(q *pop.Query, postID uuid.UUID, comments *models.Comments) error {
	// Retrieve all Comments from the DB
	return q.Where("post_id = ?", postID).Order("created_at").Eager().All(comments)
}
