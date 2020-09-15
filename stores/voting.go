package stores

import (
	"fmt"
	"rally/models"

	"github.com/gofrs/uuid"
)

type Storage interface {
	UpdateInts(f func(vals []int) error, keys ...string) ([]int, error)
	GetInt(key string, default_ *int) (int, error)
	SetInt(key string, x int) error
}

type VotingStore struct {
	s Storage
}

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
		return fmt.Errorf("internal error: incorrect payload size for VotingState")
	}
	xs[0] = state.UserBoardVotes
	xs[1] = state.UserPostVotes
	xs[2] = state.TotalPostVotes
	return nil
}

func NewVotingStore(storage Storage) VotingStore {
	return VotingStore{storage}
}

func (vs VotingStore) GetUserBoardVotes(user *models.User, board *models.Board) (int, error) {
	boardVotesKey := vs.key(user.ID, board.ID)
	votes := 0
	noVotes := 0
	var err error
	if votes, err = vs.s.GetInt(boardVotesKey, &noVotes); err != nil {
		return 0, err
	}
	return votes, nil
}

func (vs VotingStore) SetUserBoardVotes(user *models.User, board *models.Board, numVotes int) error {
	boardVotesKey := vs.key(user.ID, board.ID)
	return vs.s.SetInt(boardVotesKey, numVotes)
}

func (vs VotingStore) UpdateVotes(user *models.User, post *models.Post, f func(state *VotingState) error) error {
	boardVotesKey := vs.key(user.ID, post.BoardID)
	postVotesKey := vs.key(user.ID, post.ID)
	totalPostVotesKey := vs.key(post.ID)
	return vs.update(func(xs []int) error {
		state := &VotingState{}
		if err := state.Decode(xs); err != nil {
			return err
		}
		if err := f(state); err != nil {
			return err
		}
		if err := state.Encode(xs); err != nil {
			return err
		}
		return nil
	}, boardVotesKey, postVotesKey, totalPostVotesKey)
}

func (vs VotingStore) update(f func(xs []int) error, keys ...string) error {
	if _, err := vs.s.UpdateInts(f, keys...); err != nil {
		return err
	}
	return nil
}

func (vs VotingStore) key(IDs ...uuid.UUID) string {
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
