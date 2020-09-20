package actions

import (
	"fmt"
	"log"
	"net/http"
	"rally/models"
	"rally/services"
	"rally/stores"

	"github.com/gobuffalo/buffalo"
	"github.com/gofrs/uuid"
)

type BoardsController struct {
	AuthenticatedController

	Board        *models.Board
	RecentBoards []models.Board

	services.VotingService
	services.RecentBoardsService
	services.StarService
}

func WithBoardsController(action func(c BoardsController) error) func(c buffalo.Context) error {
	return func(ctx buffalo.Context) error {
		c := BoardsController{}
		if err := c.SetUp(ctx); err != nil {
			return err
		}
		return action(c)
	}
}

func (c *BoardsController) SetUp(ctx buffalo.Context) error {
	if err := c.AuthenticatedController.SetUp(ctx); err != nil {
		return err
	}

	boardIDParam := ctx.Param("board_id")
	if boardIDParam != "" {
		boardID, err := uuid.FromString(boardIDParam)
		if err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}
		c.Board = &models.Board{}
		if err := c.Tx.Find(c.Board, boardID); err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}

		c.Set("currentBoard", c.Board)

		c.VotingService = services.NewVotingService(stores.NewVotingStore(models.Redis),
			c.Board.VotingStrategy)

		votesRemaining, err := c.VotingService.VotesRemaining(&c.CurrentUser, c.Board)
		if err != nil {
			if err == services.ErrNoLimit {
				c.Set("isBoardVoteLimit", false)
			} else {
				log.Printf("Error determining user's remaining votes: %v", err)
			}
		}
		c.Set("isBoardVoteLimit", true)
		c.Set("votesRemaining", votesRemaining)
	}

	c.RecentBoardsService = services.NewRecentBoardsService(c.Tx)
	c.StarService = services.NewStarService(c.Tx)

	c.Set("recentBoards", func() []models.Board {
		recentBoards, err := c.RecentBoardsService.RecentBoards(&c.CurrentUser)
		if err != nil {
			return nil
		}
		return recentBoards
	})

	return nil
}

// TODO: Load board, post etc. only if the relevant RequireXYZ method called.
func (c *BoardsController) RequireBoard() error {
	if c.Board == nil {
		return c.Error(http.StatusNotFound, fmt.Errorf("missing board ID"))
	}
	return nil
}

func (c *BoardsController) RequireBoardWithWriteAccess() error {
	if err := c.RequireBoard(); err != nil {
		return err
	}
	if err := authorizeBoardManagement(c.Board, c); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	return nil
}
