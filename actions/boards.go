package actions

import (
	"net/http"
	"rally/models"

	log "github.com/sirupsen/logrus"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/x/responder"
)

// List gets all Boards. This function is mapped to the path
// GET /boards
func (c BoardsController) List() error {

	// Retrieve all Boards from the DB
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	boards := &models.Boards{}
	q := c.Tx.PaginateFromParams(c.Params())
	if err := q.All(boards); err != nil {
		return err
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		ctx.Set("pagination", q.Paginator)
		ctx.Set("boards", boards)
		return ctx.Render(http.StatusOK, r.HTML("/boards/index.plush.html"))
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.JSON(boards))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.XML(boards))
	}).Respond(c)
}

// Show gets the data for one Board. This function is mapped to
// the path GET /boards/{board_id}
func (c BoardsController) Show() error {
	if err := c.RequireBoard(); err != nil {
		return err
	}

	c.SetLastBoardID(c.Board.ID)

	// Load board's posts.
	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	posts := &models.Posts{}
	q := c.Tx.PaginateFromParams(c.Params())
	q = q.Where("(NOT draft OR (draft AND author_id = ?)) AND board_id = ?", c.CurrentUser.ID, c.Board.ID)

	order := c.Param("order")
	if order == "" {
		order = "top"
	}
	c.Set("orderClass", orderClass(order))
	if order == "newest" {
		q = q.Order("draft DESC, created_at DESC")
	} else {
		q = q.Order("draft DESC, votes DESC")
	}

	// Retrieve all Posts from the DB
	if err := q.Eager().All(posts); err != nil {
		return err
	}

	// Retrieve all Boards from the DB
	boards := &models.Boards{}
	if err := c.Tx.All(boards); err != nil {
		return err
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		ctx.Set("board", c.Board)
		ctx.Set("pagination", q.Paginator)

		ctx.Set("posts", posts)
		ctx.Set("sidebar", c.Board.Description.String != "")
		ctx.Set("boards", boards)

		return ctx.Render(http.StatusOK, r.HTML("/boards/show.plush.html"))
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.JSON(c.Board))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.XML(c.Board))
	}).Respond(c)
}

// New renders the form for creating a new Board.
// This function is mapped to the path GET /boards/new
func (c BoardsController) New() error {
	c.Set("board", models.DefaultBoard())
	return c.Render(http.StatusOK, r.HTML("/boards/new.plush.html"))
}

// Create adds a Board to the DB. This function is mapped to the
// path POST /boards
func (c BoardsController) Create() error {
	// Bind board to the html form elements
	board := &models.Board{}
	if err := c.Bind(board); err != nil {
		return err
	}

	// Validate the data from the html form
	verrs, err := c.Tx.ValidateAndCreate(board)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(ctx buffalo.Context) error {
			// Make the errors available inside the html template
			ctx.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			ctx.Set("board", board)

			return ctx.Render(http.StatusUnprocessableEntity, r.HTML("/boards/new.plush.html"))
		}).Wants("json", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	// Make the current user the owner of the board.
	member := &models.BoardMember{
		BoardID: board.ID,
		UserID:  c.CurrentUser.ID,
		IsOwner: true,
	}

	if err := c.Tx.Create(member); err != nil {
		log.Errorf("Error creating owner for board %v: %v", board.ID, err)
		return err
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a success message
		ctx.Flash().Add("success", T.Translate(ctx, "board.created.success"))

		// and redirect to the show page
		return ctx.Redirect(http.StatusSeeOther, "/boards/%v", board.ID)
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusCreated, r.JSON(board))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusCreated, r.XML(board))
	}).Respond(c)
}

// Edit renders a edit form for a Board. This function is
// mapped to the path GET /boards/{board_id}/edit
func (c BoardsController) Edit() error {
	if err := c.RequireBoardWithWriteAccess(); err != nil {
		return err
	}
	c.Set("board", c.Board)
	return c.Render(http.StatusOK, r.HTML("/boards/edit.plush.html"))
}

// Update changes a Board in the DB. This function is mapped to
// the path PUT /boards/{board_id}
func (c BoardsController) Update() error {
	if err := c.RequireBoardWithWriteAccess(); err != nil {
		return err
	}

	// Bind Board to the html form elements
	if err := c.Bind(c.Board); err != nil {
		return err
	}

	verrs, err := c.Tx.ValidateAndUpdate(c.Board)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(ctx buffalo.Context) error {
			// Make the errors available inside the html template
			ctx.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			ctx.Set("board", c.Board)

			return ctx.Render(http.StatusUnprocessableEntity, r.HTML("/boards/edit.plush.html"))
		}).Wants("json", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a success message
		ctx.Flash().Add("success", T.Translate(ctx, "board.updated.success"))

		// and redirect to the show page
		return ctx.Redirect(http.StatusSeeOther, "/boards/%v", c.Board.ID)
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JSON(c.Board))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.XML(c.Board))
	}).Respond(c)
}

// Destroy deletes a Board from the DB. This function is mapped
// to the path DELETE /boards/{board_id}
func (c BoardsController) Destroy() error {
	if err := c.RequireBoardWithWriteAccess(); err != nil {
		return err
	}

	if err := c.Tx.Destroy(c.Board); err != nil {
		return err
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a flash message
		ctx.Flash().Add("success", T.Translate(ctx, "board.destroyed.success"))

		// Redirect to the index page
		return ctx.Redirect(http.StatusSeeOther, "/boards")
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JSON(c.Board))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.XML(c.Board))
	}).Respond(c)
}
