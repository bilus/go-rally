package services

import (
	"fmt"
	"rally/models"

	"github.com/gobuffalo/pop/v5"
	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

var ErrUnauthorized = fmt.Errorf("Unauthorized")
var ErrNotFound = fmt.Errorf("Not found")

type PaginationParams interface {
	Get(key string) string
}

type PaginationResult interface {
	Paginate() string
}

type BoardsService struct {
	Tx *pop.Connection
}

func NewBoardsService(tx *pop.Connection) BoardsService {
	return BoardsService{
		Tx: tx,
	}
}

type QueryBoardsParams struct {
	User       models.User
	Pagination PaginationParams
}

type QueryBoardsResult struct {
	Boards     []models.Board
	Pagination PaginationResult
}

func (s *BoardsService) QueryBoards(params QueryBoardsParams) (QueryBoardsResult, error) {
	result := QueryBoardsResult{}
	q := s.Tx.Q()
	q = q.Where("NOT is_private OR id IN (SELECT board_id FROM board_members AS m WHERE user_id = ? AND is_owner)", params.User.ID)
	q = q.PaginateFromParams(params.Pagination)
	if err := q.All(&result.Boards); err != nil {
		return result, err
	}
	result.Pagination = q.Paginator
	return result, nil
}

type QueryBoardParams struct {
	User models.User

	BoardID                 uuid.UUID
	RequestOwnerLevelAccess bool

	IncludePosts     bool
	NewestPostsFirst bool
	PostPagination   PaginationParams
}

type QueryBoardResult struct {
	Board models.Board

	PostPagination PaginationResult
}

func (s *BoardsService) QueryBoardByID(params QueryBoardParams) (QueryBoardResult, error) {
	result := QueryBoardResult{}
	if err := s.Tx.Find(&result.Board, params.BoardID); err != nil {
		log.Errorf("Error looking for board %v: %w", params.BoardID, err)
		return result, ErrNotFound
	}
	if params.RequestOwnerLevelAccess {
		isOwner, err := s.Tx.Where("board_id = ? AND user_id = ? AND is_owner",
			params.BoardID, params.User.ID).Exists(&models.BoardMember{})
		if err != nil {
			return result, err
		}
		if !isOwner {
			return result, ErrUnauthorized
		}
	}
	if params.IncludePosts {
		q := s.Tx.PaginateFromParams(params.PostPagination)
		q = q.Where("(NOT draft OR (draft AND author_id = ?)) AND board_id = ?", params.User.ID, params.BoardID)
		if params.NewestPostsFirst {
			q = q.Order("draft DESC, created_at DESC")
		} else {
			q = q.Order("draft DESC, votes DESC")
		}

		// Retrieve all Posts from the DB
		if err := q.Eager().All(&result.Board.Posts); err != nil {
			return result, err
		}

		result.PostPagination = q.Paginator
	}
	return result, nil
}

func (s *BoardsService) DefaultBoard() models.Board {
	return *models.DefaultBoard()
}

type CreateBoardParams struct {
	User            models.User
	BoardAttributes models.Board // Use model to reduce duplication.
}

type CreateBoardResult struct {
	Board            models.Board
	ValidationErrors *validate.Errors
}

func (s *BoardsService) CreateBoard(params CreateBoardParams) (CreateBoardResult, error) {
	result := CreateBoardResult{
		Board: params.BoardAttributes,
	}
	verrs, err := s.Tx.ValidateAndCreate(&result.Board)
	if err != nil {
		return result, err
	}
	result.ValidationErrors = verrs
	if verrs.HasAny() {
		return result, nil
	}

	// Make the current user the owner of the board.
	member := models.BoardMember{
		BoardID: result.Board.ID,
		UserID:  params.User.ID,
		IsOwner: true,
	}

	if err := s.Tx.Create(&member); err != nil {
		log.Errorf("Error creating owner for board %v: %v", result.Board.ID, err)
		return result, err
	}

	return result, nil
}

type UpdateBoardParams struct {
	User    models.User
	BoardID uuid.UUID
	F       func(board *models.Board) error
}

type UpdateBoardResult struct {
	Board            models.Board
	ValidationErrors *validate.Errors
}

func (s *BoardsService) UpdateBoard(params UpdateBoardParams) (UpdateBoardResult, error) {
	result := UpdateBoardResult{}

	qr, err := s.QueryBoardByID(QueryBoardParams{
		User:                    params.User,
		BoardID:                 params.BoardID,
		RequestOwnerLevelAccess: true,
	})
	if err != nil {
		return result, err
	}

	result.Board = qr.Board

	if err := params.F(&result.Board); err != nil {
		return result, err
	}

	result.ValidationErrors, err = s.Tx.ValidateAndUpdate(&result.Board)
	return result, err
}

type DeleteBoardParams struct {
	User    models.User
	BoardID uuid.UUID
}

type DeleteBoardResult struct {
	Board models.Board
}

func (s *BoardsService) DeleteBoard(params DeleteBoardParams) (DeleteBoardResult, error) {
	result := DeleteBoardResult{}

	qr, err := s.QueryBoardByID(QueryBoardParams{
		User:                    params.User,
		BoardID:                 params.BoardID,
		RequestOwnerLevelAccess: true,
	})
	if err != nil {
		return result, err
	}

	result.Board = qr.Board

	return result, s.Tx.Destroy(&result.Board)
}
