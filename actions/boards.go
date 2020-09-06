package actions

import (
	"fmt"
	"net/http"
	"rally/models"

	log "github.com/sirupsen/logrus"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/x/responder"
	"github.com/gofrs/uuid"
)

// This file is generated by Buffalo. It offers a basic structure for
// adding, editing and deleting a page. If your model is more
// complex or you need more than the basic implementation you need to
// edit this file.

// Following naming logic is implemented in Buffalo:
// Model: Singular (Board)
// DB Table: Plural (boards)
// Resource: Plural (Boards)
// Path: Plural (/boards)
// View Template Folder: Plural (/templates/boards/)

// BoardsResource is the resource for the Board model
type BoardsResource struct {
	buffalo.Resource
}

// List gets all Boards. This function is mapped to the path
// GET /boards
func (v BoardsResource) List(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	boards := &models.Boards{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	// Retrieve all Boards from the DB
	if err := q.All(boards); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// Add the paginator to the context so it can be used in the template.
		c.Set("pagination", q.Paginator)

		c.Set("boards", boards)
		return c.Render(http.StatusOK, r.HTML("/boards/index.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(boards))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(boards))
	}).Respond(c)
}

// Show gets the data for one Board. This function is mapped to
// the path GET /boards/{board_id}
func (v BoardsResource) Show(c buffalo.Context) error {
	boardID, err := uuid.FromString(c.Param("board_id"))
	if err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Board
	board := &models.Board{}

	// To find the Board the parameter board_id is used.
	if err := tx.Find(board, boardID); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	SetLastBoardID(boardID, c)

	// Load board's posts.
	posts := &models.Posts{}

	// Paginate results. Params "page" and "per_page" control pagination.
	// Default values are "page=1" and "per_page=20".
	q := tx.PaginateFromParams(c.Params())

	currentUser, err := CurrentUser(c)
	if err != nil {
		return err
	}
	q = q.Where("(NOT draft OR (draft AND author_id = ?)) AND board_id = ?", currentUser.ID, boardID.String())

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

	boards := &models.Boards{}

	// Retrieve all Boards from the DB
	if err := tx.All(boards); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		c.Set("board", board)
		c.Set("pagination", q.Paginator)

		c.Set("posts", posts)
		c.Set("sidebar", board.Description.String != "")
		c.Set("boards", boards)

		return c.Render(http.StatusOK, r.HTML("/boards/show.plush.html"))
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(200, r.JSON(board))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(200, r.XML(board))
	}).Respond(c)
}

// New renders the form for creating a new Board.
// This function is mapped to the path GET /boards/new
func (v BoardsResource) New(c buffalo.Context) error {
	c.Set("board", models.DefaultBoard())

	return c.Render(http.StatusOK, r.HTML("/boards/new.plush.html"))
}

// Create adds a Board to the DB. This function is mapped to the
// path POST /boards
func (v BoardsResource) Create(c buffalo.Context) error {
	// Allocate an empty Board
	board := &models.Board{}

	// Bind board to the html form elements
	if err := c.Bind(board); err != nil {
		return err
	}

	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	currentUser, err := CurrentUser(c)
	if err != nil {
		return err
	}

	// Validate the data from the html form
	verrs, err := tx.ValidateAndCreate(board)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the new.html template that the user can
			// correct the input.
			c.Set("board", board)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/boards/new.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	// Make the current user the owner of the board.
	member := &models.BoardMember{
		BoardID: board.ID,
		UserID:  currentUser.ID,
		IsOwner: true,
	}

	if err := tx.Create(member); err != nil {
		log.Errorf("Error creating owner for board %v: %v", board.ID, err)
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "board.created.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/boards/%v", board.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.JSON(board))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusCreated, r.XML(board))
	}).Respond(c)
}

// Edit renders a edit form for a Board. This function is
// mapped to the path GET /boards/{board_id}/edit
func (v BoardsResource) Edit(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Board
	board := &models.Board{}

	if err := tx.Find(board, c.Param("board_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	c.Set("board", board)
	return c.Render(http.StatusOK, r.HTML("/boards/edit.plush.html"))
}

// Update changes a Board in the DB. This function is mapped to
// the path PUT /boards/{board_id}
func (v BoardsResource) Update(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Board
	board := &models.Board{}

	if err := tx.Find(board, c.Param("board_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := authorizeBoardManagement(board, c); err != nil {
		return err
	}

	// Bind Board to the html form elements
	if err := c.Bind(board); err != nil {
		return err
	}

	verrs, err := tx.ValidateAndUpdate(board)
	if err != nil {
		return err
	}

	if verrs.HasAny() {
		return responder.Wants("html", func(c buffalo.Context) error {
			// Make the errors available inside the html template
			c.Set("errors", verrs)

			// Render again the edit.html template that the user can
			// correct the input.
			c.Set("board", board)

			return c.Render(http.StatusUnprocessableEntity, r.HTML("/boards/edit.plush.html"))
		}).Wants("json", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.JSON(verrs))
		}).Wants("xml", func(c buffalo.Context) error {
			return c.Render(http.StatusUnprocessableEntity, r.XML(verrs))
		}).Respond(c)
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a success message
		c.Flash().Add("success", T.Translate(c, "board.updated.success"))

		// and redirect to the show page
		return c.Redirect(http.StatusSeeOther, "/boards/%v", board.ID)
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(board))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(board))
	}).Respond(c)
}

// Destroy deletes a Board from the DB. This function is mapped
// to the path DELETE /boards/{board_id}
func (v BoardsResource) Destroy(c buffalo.Context) error {
	// Get the DB connection from the context
	tx, ok := c.Value("tx").(*pop.Connection)
	if !ok {
		return fmt.Errorf("no transaction found")
	}

	// Allocate an empty Board
	board := &models.Board{}

	// To find the Board the parameter board_id is used.
	if err := tx.Find(board, c.Param("board_id")); err != nil {
		return c.Error(http.StatusNotFound, err)
	}

	if err := authorizeBoardManagement(board, c); err != nil {
		return err
	}

	if err := tx.Destroy(board); err != nil {
		return err
	}

	return responder.Wants("html", func(c buffalo.Context) error {
		// If there are no errors set a flash message
		c.Flash().Add("success", T.Translate(c, "board.destroyed.success"))

		// Redirect to the index page
		return c.Redirect(http.StatusSeeOther, "/boards")
	}).Wants("json", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.JSON(board))
	}).Wants("xml", func(c buffalo.Context) error {
		return c.Render(http.StatusOK, r.XML(board))
	}).Respond(c)
}
