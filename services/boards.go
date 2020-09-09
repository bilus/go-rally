package services

import (
	"rally/models"
	"sort"

	"github.com/gobuffalo/pop/v5"
)

func StarBoard(db *pop.Connection, user *models.User, board *models.Board, starred bool) error {
	user.StarBoard(board, starred)
	return db.Update(user)
}

func ListStarredBoards(db *pop.Connection, user *models.User, limit *int) ([]models.Board, error) {
	IDs := starredBoardsIDs(user)
	var q *pop.Query
	if len(IDs) > 0 {
		q = db.Where("id IN (?)", IDs...)
	} else {
		q = db.Where("FALSE")
	}
	if limit != nil {
		q = q.Limit(*limit)
	}
	boards := []models.Board{}
	if err := q.All(&boards); err != nil {
		return nil, err
	}
	return boards, nil
}

// QuickAccessBoards returns boards for the navbar selection dropdown.
func QuickAccessBoards(db *pop.Connection, user *models.User) ([]models.Board, error) {
	// TODO: Loads all boards.
	// Use Redis ZSET
	allBoards := models.Boards{}
	err := db.All(&allBoards)
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

func starredBoardsIDs(user *models.User) []interface{} {
	IDs := make([]interface{}, 0)
	for k, v := range user.StarredBoards {
		starred := v.(bool)
		if starred {
			IDs = append(IDs, k)
		}
	}
	return IDs
}
