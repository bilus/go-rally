package stores

import (
	"rally/models"
	"rally/services"
)

type VotingStore struct {
	s Storage
}

func NewVotingStore(storage Storage) VotingStore {
	return VotingStore{storage}
}

func (vs VotingStore) GetUserBoardVotes(user *models.User, board *models.Board) (int, error) {
	boardVotesKey := key(user.ID, board.ID)
	votes := 0
	noVotes := 0
	var err error
	if votes, err = vs.s.GetInt(boardVotesKey, &noVotes); err != nil {
		return 0, err
	}
	return votes, nil
}

func (vs VotingStore) SetUserBoardVotes(user *models.User, board *models.Board, numVotes int) error {
	boardVotesKey := key(user.ID, board.ID)
	return vs.s.SetInt(boardVotesKey, numVotes)
}

func (vs VotingStore) UpdateVotes(user *models.User, post *models.Post, f func(state *services.VotingState) error) error {
	boardVotesKey := key(user.ID, post.BoardID)
	postVotesKey := key(user.ID, post.ID)
	totalPostVotesKey := key(post.ID)
	return vs.update(func(xs []int) error {
		state := &services.VotingState{}
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

// TODO: Cleaning up redis keys when posts/boards deleted.
