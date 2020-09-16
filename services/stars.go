package services

import (
	"rally/models"

	"github.com/gobuffalo/pop/v5"
)

type StarService struct {
	db *pop.Connection
}

func NewStarService(db *pop.Connection) StarService {
	return StarService{db}
}

func (s StarService) StarBoard(user *models.User, board *models.Board, starred bool) error {
	user.StarBoard(board, starred)
	return s.db.Update(user)
}

func (s StarService) ListStarredBoards(user *models.User, limit *int) ([]models.Board, error) {
	IDs := starredBoardsIDs(user)
	var q *pop.Query
	if len(IDs) > 0 {
		q = s.db.Where("id IN (?)", IDs...)
	} else {
		q = s.db.Where("FALSE")
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
