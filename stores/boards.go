package stores

import (
	"database/sql"
	"rally/models"
	"rally/services"

	"github.com/gobuffalo/pop/v5"
	"github.com/gofrs/uuid"
)

type BoardsStore struct {
	tx *pop.Connection
}

func NewBoardsStore(tx *pop.Connection) BoardsStore {
	return BoardsStore{tx}
}

// TODO: Type assertion.

func (s BoardsStore) ListBoards(user models.User, pg services.PaginationParams) ([]models.Board, services.PaginationResult, error) {
	boards := []models.Board{}
	q := s.tx.Q()
	q = q.Where("NOT is_private OR id IN (SELECT board_id FROM board_members AS m WHERE user_id = ? AND is_owner)", user.ID)
	q = q.PaginateFromParams(pg)
	if err := q.All(&boards); err != nil {
		return nil, nil, err
	}
	return boards, q.Paginator, nil
}

func (s BoardsStore) FindBoardByID(boardID uuid.UUID) (board models.Board, found bool, err error) {
	if err = s.tx.Find(&board, boardID); err != nil {
		if err == sql.ErrNoRows {
			err = nil
		}
		found = false
	} else {
		found = true
	}
	return
}

func (s BoardsStore) IsOwner(boardID, userID uuid.UUID) (bool, error) {
	return s.tx.Where("board_id = ? AND user_id = ? AND is_owner",
		boardID, userID).Exists(&models.BoardMember{})
}

func (s BoardsStore) ListBoardPosts(boardID, userID uuid.UUID, newestFirst, archived bool, pg services.PaginationParams) ([]models.Post, services.PaginationResult, error) {
	posts := []models.Post{}
	q := s.tx.PaginateFromParams(pg)
	q = q.Where("(NOT draft OR (draft AND author_id = ?)) AND board_id = ?", userID, boardID)
	if archived {
		q = q.Where("archived")
	} else {
		q = q.Where("NOT archived")
	}
	if newestFirst {
		q = q.Order("draft DESC, created_at DESC")
	} else {
		q = q.Order("draft DESC, votes DESC")
	}

	// Retrieve all Posts from the DB
	if err := q.Eager().All(&posts); err != nil {
		return nil, nil, err
	}
	return posts, q.Paginator, nil
}

func (s BoardsStore) CreateBoard(board *models.Board) error {
	return s.tx.Create(board)
}

func (s BoardsStore) AddBoardOwner(boardID, userID uuid.UUID) error {
	// Make the current user the owner of the board.
	member := models.BoardMember{
		BoardID: boardID,
		UserID:  userID,
		IsOwner: true,
	}

	return s.tx.Create(&member)
}

func (s BoardsStore) SaveBoard(board models.Board) error {
	return s.tx.Update(&board)
}

func (s BoardsStore) DeleteBoard(board models.Board) error {
	return s.tx.Destroy(&board)
}
