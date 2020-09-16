package services

import (
	"rally/models"
	"sort"

	"github.com/gobuffalo/pop/v5"
)

type RecentBoardsService struct {
	db *pop.Connection
}

func NewRecentBoardsService(db *pop.Connection) RecentBoardsService {
	return RecentBoardsService{db}
}

// RecentBoards returns boards for the navbar selection dropdown.
func (s RecentBoardsService) RecentBoards(user *models.User) ([]models.Board, error) {
	// TODO: Loads all boards.
	// Use Redis ZSET
	allBoards := models.Boards{}
	err := s.db.All(&allBoards)
	if err != nil {
		return nil, err
	}

	allBoards.AddUserContext(user)

	allBoardsSlice := allBoards.Slice()
	sort.Slice(allBoards, func(i, j int) bool {
		return priority(&allBoardsSlice[i]) > priority(&allBoardsSlice[j])
	})

	maxBoards := 5
	boards := make([]models.Board, 0, maxBoards)
	for i, b := range allBoards {
		if i >= maxBoards {
			break
		}
		boards = append(boards, b)
	}
	return boards, err
}

func priority(board *models.Board) int {
	if board.UserStarred {
		return 1
	} else {
		return 0
	}
}
