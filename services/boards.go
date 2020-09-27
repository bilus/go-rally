package services

import (
	"fmt"
	"rally/models"

	"github.com/gobuffalo/validate/v3"
	"github.com/gofrs/uuid"
)

var ErrUnauthorized = fmt.Errorf("Unauthorized")
var ErrNotFound = fmt.Errorf("Not found")

type PaginationParams interface {
	Get(key string) string
}

type PaginationResult interface {
	Paginate() string
}

type BoardsStore interface {
	ListBoards(user models.User, pg PaginationParams) ([]models.Board, PaginationResult, error)
	FindBoardByID(boardID uuid.UUID) (models.Board, bool, error)
	IsOwner(boardID, userID uuid.UUID) (bool, error)
	ListBoardPosts(boardID, userID uuid.UUID, newestFirst bool, pg PaginationParams) ([]models.Post, PaginationResult, error)
	CreateBoard(board *models.Board) error
	AddBoardOwner(boardID, userID uuid.UUID) error
	SaveBoard(board models.Board) error
	DeleteBoard(board models.Board) error
}

type BoardsService struct {
	store            BoardsStore
	reactionsService ReactionsService
}

func NewBoardsService(store BoardsStore, reactionsService ReactionsService) BoardsService {
	return BoardsService{
		store:            store,
		reactionsService: reactionsService,
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

func (s *BoardsService) QueryBoards(params QueryBoardsParams) (result QueryBoardsResult, err error) {
	result.Boards, result.Pagination, err = s.store.ListBoards(params.User, params.Pagination)
	return
}

type QueryBoardParams struct {
	User models.User

	BoardID                 uuid.UUID
	RequestOwnerLevelAccess bool

	IncludePosts     bool
	NewestPostsFirst bool
	IncludeReactions bool
	PostPagination   PaginationParams
}

type QueryBoardResult struct {
	Board models.Board

	PostPagination PaginationResult
}

func (s *BoardsService) QueryBoardByID(params QueryBoardParams) (result QueryBoardResult, err error) {
	var found bool
	result.Board, found, err = s.store.FindBoardByID(params.BoardID)
	if err != nil {
		return
	}
	if !found {
		err = ErrNotFound
		return
	}
	result.Board.UserIsOwner, err = s.store.IsOwner(params.BoardID,
		params.User.ID)
	if err != nil {
		return
	}
	if params.RequestOwnerLevelAccess {
		if !result.Board.UserIsOwner {
			err = ErrUnauthorized
			return
		}
	}
	if params.IncludePosts {
		result.Board.Posts, result.PostPagination, err = s.store.ListBoardPosts(params.BoardID, params.User.ID, params.NewestPostsFirst, params.PostPagination)
		if err != nil {
			return
		}
		if params.IncludeReactions {
			for i, post := range result.Board.Posts {
				result.Board.Posts[i].Reactions, err = s.reactionsService.ListReactionsToPost(&params.User, &post)
				if err != nil {
					return
				}
			}
		}
	}
	return
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

func (s *BoardsService) CreateBoard(params CreateBoardParams) (result CreateBoardResult, err error) {
	result.Board = params.BoardAttributes
	result.ValidationErrors, err = result.Board.ValidateCreate(nil) // TODO
	if err != nil || result.ValidationErrors.HasAny() {
		return
	}
	err = s.store.CreateBoard(&result.Board)
	if err != nil {
		return
	}
	err = s.store.AddBoardOwner(result.Board.ID, params.User.ID)
	return
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

func (s *BoardsService) UpdateBoard(params UpdateBoardParams) (result UpdateBoardResult, err error) {
	var qr QueryBoardResult
	qr, err = s.QueryBoardByID(QueryBoardParams{
		User:                    params.User,
		BoardID:                 params.BoardID,
		RequestOwnerLevelAccess: true,
	})
	if err != nil {
		return
	}

	result.Board = qr.Board

	if err = params.F(&result.Board); err != nil {
		return
	}

	result.ValidationErrors, err = result.Board.ValidateUpdate(nil)
	if err != nil || result.ValidationErrors.HasAny() {
		return
	}
	err = s.store.SaveBoard(result.Board)
	return
}

type DeleteBoardParams struct {
	User    models.User
	BoardID uuid.UUID
}

type DeleteBoardResult struct {
	Board models.Board
}

func (s *BoardsService) DeleteBoard(params DeleteBoardParams) (result DeleteBoardResult, err error) {
	var qr QueryBoardResult
	qr, err = s.QueryBoardByID(QueryBoardParams{
		User:                    params.User,
		BoardID:                 params.BoardID,
		RequestOwnerLevelAccess: true,
	})
	if err != nil {
		return
	}
	err = s.store.DeleteBoard(result.Board)
	if err != nil {
		return
	}

	result.Board = qr.Board
	return
}
