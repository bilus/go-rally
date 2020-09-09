package models

import (
	"fmt"

	"github.com/gobuffalo/nulls"
	"github.com/gofrs/uuid"

	log "github.com/sirupsen/logrus"
)

var ErrLimit = fmt.Errorf("limit reached")
var ErrNoLimit = fmt.Errorf("board has no vote limit")

type Store interface {
	UpdateInts(f func(vals []int) error, keys ...string) ([]int, error)
	GetInt(key string, default_ *int) (int, error)
	SetInt(key string, x int) error
}

type VotingStrategy struct {
	BoardMax nulls.Int `json:"board_max"`
}

func (s VotingStrategy) VotesRemaining(store Store, user *User, board *Board) (int, error) {
	if !s.BoardMax.Valid {
		return 0, ErrNoLimit
	}
	boardVotesKey := s.key(user.ID, board.ID)
	noVotes := 0
	votes, err := store.GetInt(boardVotesKey, &noVotes)
	if err != nil {
		return 0, err
	}

	return s.BoardMax.Int - votes, nil
}

func (s VotingStrategy) Downvote(store Store, user *User, post *Post) (totalPostVotes int, success bool, err error) {
	return s.update(store, user, post, func(xs []int) error {
		boardVotes := xs[0]
		if boardVotes <= 0 {
			return ErrLimit
		}
		xs[0] = boardVotes - 1
		postVotes := xs[1]
		if postVotes <= 0 {
			return ErrLimit
		}
		xs[1] = postVotes - 1
		totalPostVotes := xs[2]
		xs[2] = totalPostVotes - 1 // TODO: No negativity check?
		return nil
	})
}

func (s VotingStrategy) Upvote(store Store, user *User, post *Post) (totalPostVotes int, success bool, err error) {
	return s.update(store, user, post, func(xs []int) error {
		boardVotes := xs[0]
		if s.BoardMax.Valid {
			if boardVotes >= s.BoardMax.Int {
				return ErrLimit
			}
		}
		xs[0] = boardVotes + 1
		postVotes := xs[1]
		xs[1] = postVotes + 1
		totalPostVotes := xs[2]
		xs[2] = totalPostVotes + 1
		return nil
	})
}

func (s VotingStrategy) Refill(store Store, board *Board, userIDs ...uuid.UUID) error {
	var lastErr error
	for _, ID := range userIDs {
		boardVotesKey := s.key(ID, board.ID)
		if err := store.SetInt(boardVotesKey, 0); err != nil {
			log.Errorf("Failed to refill board %v votes for user %v: %v", board.ID, ID, err)
			lastErr = err

		}
	}
	return lastErr
}

func (s VotingStrategy) update(store Store, user *User, post *Post, f func(xs []int) error) (postVotes int, success bool, err error) {
	boardVotesKey := s.key(user.ID, post.BoardID)
	postVotesKey := s.key(user.ID, post.ID)
	totalPostVotesKey := s.key(post.ID)

	xs, err := store.UpdateInts(f, boardVotesKey, postVotesKey, totalPostVotesKey)

	if err == ErrLimit {
		zero := 0
		votes, err := store.GetInt(totalPostVotesKey, &zero)
		if err != nil {
			return 0, false, err
		}
		return votes, false, nil
	}
	if err != nil {
		return 0, false, err
	}
	return xs[2], true, nil
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
