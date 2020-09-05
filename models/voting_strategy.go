package models

import (
	"fmt"

	"github.com/gofrs/uuid"
)

type Store interface {
	UpdateInts(f func(vals []int) error, keys ...string) error
	GetInt(key string, default_ *int) (int, error)
}

type VotingStrategy struct {
	BoardMax int `json:"board_max"`
}

func (s VotingStrategy) VotesRemaining(store Store, user *User, board *Board) (int, error) {
	k := s.key(user.ID, board.ID)
	noVotes := 0
	votes, err := store.GetInt(k, &noVotes)
	if err != nil {
		return 0, err
	}

	return s.BoardMax - votes, nil
}

func (s VotingStrategy) Downvote(store Store, user *User, post *Post) (bool, error) {
	boardKey := s.key(user.ID, post.BoardID)
	pageKey := s.key(user.ID, post.ID)
	ErrLimit := fmt.Errorf("limit reached")
	err := store.UpdateInts(func(xs []int) error {
		boardVotes := xs[0]
		if boardVotes <= 0 {
			return ErrLimit
		}
		xs[0] = boardVotes - 1
		pageVotes := xs[1]
		if pageVotes <= 0 {
			return ErrLimit
		}
		xs[1] = pageVotes - 1
		return nil
	}, boardKey, pageKey)
	if err == ErrLimit {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s VotingStrategy) Upvote(store Store, user *User, post *Post) (bool, error) {
	boardKey := s.key(user.ID, post.BoardID)
	pageKey := s.key(user.ID, post.ID)
	ErrLimit := fmt.Errorf("limit reached")
	err := store.UpdateInts(func(xs []int) error {
		boardVotes := xs[0]
		if boardVotes >= s.BoardMax {
			return ErrLimit
		}
		xs[0] = boardVotes + 1
		pageVotes := xs[1]
		xs[1] = pageVotes + 1
		return nil
	}, boardKey, pageKey)
	if err == ErrLimit {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (s VotingStrategy) key(IDs ...uuid.UUID) string {
	key := ""
	for _, ID := range IDs {
		if key != "" {
			key += ":"
		}
		key += ID.String()
	}
	return key
}

// TODO: Cleaning up redis keys when posts/boards deleted.
