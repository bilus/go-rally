package services

import (
	"fmt"
	"rally/models"

	log "github.com/sirupsen/logrus"
)

var ErrLimit = fmt.Errorf("limit reached")
var ErrNoLimit = fmt.Errorf("board has no vote limit")

type VotingState struct {
	UserBoardVotes int
	UserPostVotes  int
	TotalPostVotes int
}

func (state *VotingState) Decode(xs []int) error {
	if len(xs) != 3 {
		return fmt.Errorf("internal error: incorrect payload size for VotingState")
	}
	state.UserBoardVotes = xs[0]
	state.UserPostVotes = xs[1]
	state.TotalPostVotes = xs[2]
	return nil
}

func (state *VotingState) Encode(xs []int) error {
	if len(xs) != 3 {
		return fmt.Errorf("internal error: incorrect payload size for services.VotingState")
	}
	xs[0] = state.UserBoardVotes
	xs[1] = state.UserPostVotes
	xs[2] = state.TotalPostVotes
	return nil
}

type Store interface {
	GetUserBoardVotes(user *models.User, board *models.Board) (int, error)
	SetUserBoardVotes(user *models.User, board *models.Board, numVotes int) error
	UpdateVotes(user *models.User, post *models.Post, f func(state *VotingState) error) error
}

type VotingService struct {
	strategy models.VotingStrategy
	store    Store
}

func NewVotingService(store Store, strategy models.VotingStrategy) VotingService {
	return VotingService{
		strategy: strategy,
		store:    store,
	}
}

func (s VotingService) VotesRemaining(user *models.User, board *models.Board) (int, error) {
	if !s.strategy.BoardMax.Valid {
		return 0, ErrNoLimit
	}
	numVotes, err := s.store.GetUserBoardVotes(user, board)
	if err != nil {
		return 0, err
	}
	numVotesRemain := s.strategy.BoardMax.Int - numVotes
	if numVotesRemain < 0 {
		log.Debugf("Negative amount of votes for user %v: %v", user.ID, numVotesRemain)
		return 0, nil
	}
	return numVotesRemain, nil
}

func (s VotingService) Downvote(user *models.User, post *models.Post) (int, bool, error) {
	var totalPostVotes int
	err := s.store.UpdateVotes(user, post, func(state *VotingState) error {
		if state.UserBoardVotes <= 0 {
			return ErrLimit
		}
		state.UserBoardVotes--
		if state.UserPostVotes <= 0 {
			return ErrLimit
		}
		state.UserPostVotes--
		state.TotalPostVotes--
		totalPostVotes = state.TotalPostVotes // Return value.
		return nil
	})
	if err != nil {
		if err == ErrLimit {
			return totalPostVotes, false, nil
		} else {
			return 0, false, err
		}
	}
	return totalPostVotes, true, nil
}

func (s VotingService) Upvote(user *models.User, post *models.Post) (int, bool, error) {
	var totalPostVotes int
	err := s.store.UpdateVotes(user, post, func(state *VotingState) error {
		if s.strategy.BoardMax.Valid {
			if state.UserBoardVotes >= s.strategy.BoardMax.Int {
				totalPostVotes = state.TotalPostVotes // Return value.
				return ErrLimit
			}
		}
		state.UserBoardVotes++
		state.UserPostVotes++
		state.TotalPostVotes++
		totalPostVotes = state.TotalPostVotes // Return value.
		return nil
	})
	if err != nil {
		if err == ErrLimit {
			return totalPostVotes, false, nil
		} else {
			return 0, false, err
		}
	}
	return totalPostVotes, true, nil
}

func (s VotingService) Refill(board *models.Board, users ...models.User) error {
	var lastErr error
	for _, user := range users {
		if err := s.store.SetUserBoardVotes(&user, board, 0); err != nil {
			lastErr = err
		}
	}
	return lastErr
}
