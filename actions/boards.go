package actions

import (
	"net/http"
	"rally/models"
	"rally/services"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/x/responder"
)

// List gets all Boards. This function is mapped to the path
// GET /boards
func (c BoardsController) List() error {
	result, err := c.BoardsService.QueryBoards(
		services.QueryBoardsParams{
			User:       c.CurrentUser,
			Pagination: c.PaginationParams,
		})
	if err != nil {
		return err
	}
	return responder.Wants("html", func(ctx buffalo.Context) error {
		ctx.Set("pagination", result.Pagination)
		ctx.Set("boards", result.Boards)
		return ctx.Render(http.StatusOK, r.HTML("/boards/index.plush.html"))
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.JSON(result.Boards))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.XML(result.Boards))
	}).Respond(c)
}

// Show gets the data for one Board. This function is mapped to
// the path GET /boards/{board_id}
func (c BoardsController) Show() error {
	c.SetLastBoardID(c.Board.ID)
	result, err := c.BoardsService.QueryBoardByID(
		services.QueryBoardParams{
			User:             c.CurrentUser,
			BoardID:          c.Board.ID,
			IncludePosts:     true,
			NewestPostsFirst: c.Param("order") == "newest",
			PostPagination:   c.PaginationParams,
		})
	if err != nil {
		return err
	}
	c.Set("orderClass", orderClass(c.Param("order")))
	return responder.Wants("html", func(ctx buffalo.Context) error {
		ctx.Set("board", result.Board)
		ctx.Set("pagination", result.PostPagination)

		ctx.Set("posts", result.Board.Posts)
		ctx.Set("sidebar", result.Board.Description.String != "")

		return ctx.Render(http.StatusOK, r.HTML("/boards/show.plush.html"))
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.JSON(result.Board))
	}).Wants("xml", func(ctx buffalo.Context) error {
		return ctx.Render(200, r.XML(result.Board))
	}).Respond(c)
}

// New renders the form for creating a new Board.
// This function is mapped to the path GET /boards/new
func (c BoardsController) New() error {
	c.Set("board", c.BoardsService.DefaultBoard())
	return c.Render(http.StatusOK, r.HTML("/boards/new.plush.html"))
}

// Create adds a Board to the DB. This function is mapped to the
// path POST /boards
func (c BoardsController) Create() error {
	// Bind board to the html form elements
	params := services.CreateBoardParams{}
	if err := c.Bind(&params.BoardAttributes); err != nil {
		return err
	}

	result, err := c.BoardsService.CreateBoard(params)
	if err != nil {
		return err
	}

	if result.ValidationErrors.HasAny() {
		return responder.Wants("html", func(ctx buffalo.Context) error {
			// Make the errors available inside the html template
			ctx.Set("errors", result.ValidationErrors)

			// Render again the new.html template that the user can
			// correct the input.
			ctx.Set("board", result.Board)

			return ctx.Render(http.StatusUnprocessableEntity, r.HTML("/boards/new.plush.html"))
		}).Wants("json", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.JSON(result.ValidationErrors))
		}).Respond(c)
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a success message
		ctx.Flash().Add("success", T.Translate(ctx, "board.created.success"))

		// and redirect to the show page
		return ctx.Redirect(http.StatusSeeOther, "/boards/%v", result.Board.ID)
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusCreated, r.JSON(result.Board))
	}).Respond(c)
}

// Edit renders a edit form for a Board. This function is
// mapped to the path GET /boards/{board_id}/edit
func (c BoardsController) Edit() error {
	result, err := c.BoardsService.QueryBoardByID(services.QueryBoardParams{
		User:                    c.CurrentUser,
		BoardID:                 c.Board.ID,
		RequestOwnerLevelAccess: true,
	})
	if err != nil {
		return err
	}
	c.Set("board", result.Board)
	return c.Render(http.StatusOK, r.HTML("/boards/edit.plush.html"))
}

// Update changes a Board in the DB. This function is mapped to
// the path PUT /boards/{board_id}
func (c BoardsController) Update() error {
	result, err := c.BoardsService.UpdateBoard(
		services.UpdateBoardParams{
			User:    c.CurrentUser,
			BoardID: c.Board.ID,
			F:       func(b *models.Board) error { return c.Bind(b) },
		})

	if err != nil {
		return err
	}

	if result.ValidationErrors.HasAny() {
		return responder.Wants("html", func(ctx buffalo.Context) error {
			// Make the errors available inside the html template
			ctx.Set("errors", result.ValidationErrors)

			// Render again the edit.html template that the user can
			// correct the input.
			ctx.Set("board", result.Board)

			return ctx.Render(http.StatusUnprocessableEntity, r.HTML("/boards/edit.plush.html"))
		}).Wants("json", func(ctx buffalo.Context) error {
			return ctx.Render(http.StatusUnprocessableEntity, r.JSON(result.ValidationErrors))
		}).Respond(c)
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a success message
		ctx.Flash().Add("success", T.Translate(ctx, "board.updated.success"))

		// and redirect to the show page
		return ctx.Redirect(http.StatusSeeOther, "/boards/%v", result.Board.ID)
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JSON(result.Board))
	}).Respond(c)
}

// Destroy deletes a Board from the DB. This function is mapped
// to the path DELETE /boards/{board_id}
func (c BoardsController) Destroy() error {
	result, err := c.BoardsService.DeleteBoard(
		services.DeleteBoardParams{
			User:    c.CurrentUser,
			BoardID: c.Board.ID,
		})
	if err != nil {
		return err
	}

	return responder.Wants("html", func(ctx buffalo.Context) error {
		// If there are no errors set a flash message
		ctx.Flash().Add("success", T.Translate(ctx, "board.destroyed.success"))

		// Redirect to the index page
		return ctx.Redirect(http.StatusSeeOther, "/dashboard")
	}).Wants("json", func(ctx buffalo.Context) error {
		return ctx.Render(http.StatusOK, r.JSON(result.Board))
	}).Respond(c)
}
