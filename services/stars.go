package services

import (
	"rally/models"

	"github.com/gobuffalo/pop/v5"
)

func StarBoard(db *pop.Connection, user *models.User, board *models.Board, starred bool) error {
	user.StarBoard(board, starred)
	return db.Update(user)
}
