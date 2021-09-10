package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/x/responder"
)

// Show gets the data for one Post.
func (c PostsController) Show() error {
	if err := c.RequirePost(); err != nil {
		return err
	}

	reactions, err := c.ReactionsService.ListReactionsToPost(&c.CurrentUser, c.Post)
	if err != nil {
		return err
	}
	c.Post.Reactions = reactions

	return responder.Wants("html", func(ctx buffalo.Context) error {
		ctx.Set("post", c.Post)
		ctx.Set("board", c.Board)
		return ctx.Render(http.StatusOK, r.HTML("/posts/show.plush.html"))
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.JSON(c.Post))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.XML(c.Post))
	}).Respond(c)
}

// Create adds a Post to the DB.
func (c PostsController) Create() error {
	if err := c.RequireBoard(); err != nil {
		return err
	}

	// Allocate an empty Post draft.
	c.Post = &models.Post{
		Draft:    true,
		BoardID:  c.Board.ID,
		AuthorID: c.CurrentUser.ID,
		Author:   c.CurrentUser,
	}

	// Validate the data from the html form
	verrs, err := c.Tx.ValidateAndCreate(c.Post)
	if err != nil {
		return err
	}
	if verrs.HasAny() {
		return fmt.Errorf("validation failed when creating an empty new draft: %q", verrs.String())
	}

	c.Set("post", c.Post)
	return c.Render(http.StatusOK, r.JavaScript("/posts/create.plush.js"))
}

// Edit renders a edit form for a Post.
func (c PostsController) Edit() error {
	if err := c.RequirePostWithWriteAccess(); err != nil {
		return err
	}
	c.Set("post", c.Post)
	c.Set("board", c.Board)
	return c.Render(http.StatusOK, r.HTML("/posts/edit.plush.html"))
}

// Update changes a Post in the DB.
func (c PostsController) Update() error {
	if err := c.RequirePostWithWriteAccess(); err != nil {
		return err
	}

	// Bind Post to the html form elements and update it.
	if err := c.Bind(c.Post); err != nil {
		return err
	}
	if c.Param("archived") == "true" {
		c.Post.Archived = true
	} else if c.Param("archived") == "false" {
		c.Post.Archived = false
	}

	verrs, err := c.Tx.ValidateAndUpdate(c.Post)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(ctx buffalo.Context) error {
			// Make the errors available inside the html template
			ctx.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			ctx.Set("post", c.Post)
			ctx.Set("board", c.Board)

			return ctx.Render(http.StatusUnprocessableEntity, r.HTML("/posts/edit.plush.html"))
		}).Wants("json", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a success message
		ctx.Flash().Add("success", T.Translate(c, "post.updated.success"))

		// and redirect to the show page
		return ctx.Redirect(http.StatusSeeOther, "/boards/%v/posts/%v", c.Post.BoardID, c.Post.ID) // TODO: Use path helper.
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JSON(c.Post))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.XML(c.Post))
	}).Respond(c)
}

// Destroy deletes a Post from the DB.
func (c PostsController) Destroy() error {
	if err := c.RequirePostWithWriteAccess(); err != nil {
		return err
	}

	if err := c.Tx.Destroy(c.Post); err != nil {
		return err
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a flash message
		ctx.Flash().Add("success", T.Translate(c, "post.destroyed.success"))

		// Redirect to the index page
		return ctx.Redirect(http.StatusSeeOther, "/boards/%v", c.Post.BoardID)
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JSON(c.Post))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.XML(c.Post))
	}).Respond(c)
}
