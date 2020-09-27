package actions

import (
	"fmt"
	"net/http"
	"rally/models"
	"rally/services"
	"rally/stores"

	log "github.com/sirupsen/logrus"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/gobuffalo/plush/v4"
	"github.com/gofrs/uuid"
)

type BoardsController struct {
	AuthenticatedController

	Board        *models.Board
	RecentBoards []models.Board

	BoardsService       services.BoardsService
	VotingService       services.VotingService
	RecentBoardsService services.RecentBoardsService
	StarService         services.StarService
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

	c.BoardsService = services.NewBoardsService(stores.NewBoardsStore(c.Tx),
		// TODO: Duplicated in PostsController.
		services.NewReactionsService(stores.NewReactionStore(models.Redis), c.Tx))

	boardIDParam := ctx.Param("board_id")
	if boardIDParam != "" {
		boardID, err := uuid.FromString(boardIDParam)
		if err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}
		// TODO: It's only here because of actions outside boards.go
		result, err := c.BoardsService.QueryBoardByID(services.QueryBoardParams{
			User:    c.CurrentUser,
			BoardID: boardID,
		})
		if err != nil {
			return ctx.Error(http.StatusNotFound, err)
		}
		c.Board = &result.Board
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
			log.Warnf("Error in recentBoards helper: %v", err)
			return nil
		}
		return recentBoards
	})

	c.RegisterHelpers(render.Helpers{
		"canManagePost": func(post interface{}) bool {
			p := toPostPtr(post)
			return c.authorizePostManagement(p) == nil
		},
		"canManageBoard": func(board interface{}) bool {
			b := toBoardPtr(board)
			return c.authorizeBoardManagement(b) == nil
		},
		"isBoardStarred": func(board interface{}, help plush.HelperContext) bool {
			b := toBoardPtr(board)
			u, err := CurrentUser(help)
			if err != nil {
				log.Errorf("Unable to retrieve current user to determine if board starred: %v", err)
				return false
			}
			return u.IsBoardStarred(b)
		},
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
	if err := c.authorizeBoardManagement(c.Board); err != nil {
		return c.Error(http.StatusUnauthorized, err)
	}

	return nil
}

func (c BoardsController) authorizeBoardManagement(board *models.Board) error {
	if !board.UserIsOwner {
		return fmt.Errorf("no permission to manage board %q", board.ID)
	}
	return nil
}

func (c BoardsController) authorizePostManagement(post *models.Post) error {
	if c.authorizeBoardManagement(c.Board) == nil {
		return nil
	}
	if post.AuthorID == c.CurrentUser.ID {
		return nil
	}
	return fmt.Errorf("no permission to manage post %q", post.ID)
}
