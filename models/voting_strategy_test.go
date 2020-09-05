package models_test

import (
	"fmt"
	"rally/models"
)

type FakeStore struct {
	m map[string]int
}

func NewFakeStore() FakeStore {
	return FakeStore{make(map[string]int)}
}

func (s FakeStore) GetInt(key string, default_ *int) (int, error) {
	i, ok := s.m[key]
	if !ok {
		if default_ != nil {
			return *default_, nil
		}
		return 0, fmt.Errorf("missing value at %q and no default", key)
	}
	return i, nil
}

func (s FakeStore) UpdateInts(f func(vals []int) error, keys ...string) error {
	vals := make([]int, len(keys))
	for i, k := range keys {
		v, ok := s.m[k]
		if !ok {
			v = 0
		}
		vals[i] = v
	}
	if err := f(vals); err != nil {
		return err
	}
	for i, k := range keys {
		s.m[k] = vals[i]
	}
	return nil
}

// Can upvote up to a limit.
func (t *ModelSuite) Test_VotingStrategy_UpvotingUpperLimit() {
	store := NewFakeStore()

	u := t.MustCreateUser()
	b := t.MustCreateBoardWithVotingStrategy(models.VotingStrategy{BoardMax: 1})

	p := t.MustCreatePost(t.ValidPost(b, u))

	s := b.VotingStrategy
	count, err := s.VotesRemaining(store, u, b)
	t.NoError(err)
	t.Equal(1, count)

	upvoted, err := s.Upvote(store, u, p)
	t.NoError(err)
	t.True(upvoted)

	count, err = s.VotesRemaining(store, u, b)
	t.NoError(err)
	t.Equal(0, count)

	upvoted, err = s.Upvote(store, u, p)
	t.NoError(err)
	t.False(upvoted)
}

// Boards keep independent vote counts.
func (t *ModelSuite) Test_VotingStrategy_BoardIsolation() {
	store := NewFakeStore()

	u := t.MustCreateUser()

	b1 := t.MustCreateBoardWithVotingStrategy(models.VotingStrategy{BoardMax: 1})
	b2 := t.MustCreateBoardWithVotingStrategy(models.VotingStrategy{BoardMax: 1})

	p1 := t.MustCreatePost(t.ValidPost(b1, u))
	p2 := t.MustCreatePost(t.ValidPost(b2, u))

	_, err := b1.Upvote(store, u, p1)
	t.NoError(err)

	count, err := b1.VotesRemaining(store, u, b1)
	t.NoError(err)
	t.Equal(0, count)

	count, err = b2.VotesRemaining(store, u, b2)
	t.NoError(err)
	t.Equal(1, count)
}

// Can take votes back until going to zero.
func (t *ModelSuite) Test_VotingStrategy_Backsies() {
	store := NewFakeStore()

	u := t.MustCreateUser()
	b := t.MustCreateBoardWithVotingStrategy(models.VotingStrategy{BoardMax: 1})

	p := t.MustCreatePost(t.ValidPost(b, u))

	s := b.VotingStrategy

	upvoted, err := s.Upvote(store, u, p)
	t.NoError(err)
	t.True(upvoted)

	count, err := s.VotesRemaining(store, u, b)
	t.NoError(err)
	t.Equal(0, count)

	downvoted, err := s.Downvote(store, u, p)
	t.NoError(err)
	t.True(downvoted)

	count, err = s.VotesRemaining(store, u, b)
	t.NoError(err)
	t.Equal(1, count)

	downvoted, err = s.Downvote(store, u, p)
	t.NoError(err)
	t.False(downvoted)

	count, err = s.VotesRemaining(store, u, b)
	t.NoError(err)
	t.Equal(1, count)

	upvoted, err = s.Upvote(store, u, p)
	t.NoError(err)
	t.True(upvoted)

	count, err = s.VotesRemaining(store, u, b)
	t.NoError(err)
	t.Equal(0, count)

}

// Cannot take votes back for non-upvoted posts.
func (t *ModelSuite) Test_VotingStrategy_BacksiesForUpvotedPosts() {
	store := NewFakeStore()

	u := t.MustCreateUser()

	b := t.MustCreateBoardWithVotingStrategy(models.VotingStrategy{BoardMax: 1})

	p1 := t.MustCreatePost(t.ValidPost(b, u))
	p2 := t.MustCreatePost(t.ValidPost(b, u))

	_, err := b.Upvote(store, u, p1)
	t.NoError(err)

	count, err := b.VotesRemaining(store, u, b1)
	t.NoError(err)
	t.Equal(0, count)

	count, err = b.VotesRemaining(store, u, b2)
	t.NoError(err)
	t.Equal(0, count)

	downvoted, err := b.Downvote(store, u, p2)
	t.NoError(err)
	t.False(downvoted)

	count, err = b.VotesRemaining(store, u, b1)
	t.NoError(err)
	t.Equal(0, count)

	count, err = b.VotesRemaining(store, u, b2)
	t.NoError(err)
	t.Equal(0, count)

	downvoted, err = b.Downvote(store, u, p1)
	t.NoError(err)
	t.True(downvoted)

	count, err = b.VotesRemaining(store, u, b1)
	t.NoError(err)
	t.Equal(1, count)

	count, err = b.VotesRemaining(store, u, b2)
	t.NoError(err)
	t.Equal(1, count)
}
