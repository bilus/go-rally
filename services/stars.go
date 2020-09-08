package services

import (
	"rally/models"

	"github.com/gobuffalo/pop/v5"
)

func StarBoard(db *pop.Connection, user *models.User, board *models.Board, starred bool) error {
	user.StarBoard(board, starred)
	return db.Update(user)
}

func ListStarredBoards(db *pop.Connection, user *models.User) (*pop.Query, error) {
	IDs := make([]interface{}, 0)
	for k, v := range user.StarredBoards {
		starred := v.(bool)
		if starred {
			IDs = append(IDs, k)
		}
	}

	if len(IDs) > 0 {
		return db.Where("id IN (?)", IDs...), nil
	} else {
		return db.Where("FALSE"), nil
	}
}
