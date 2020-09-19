package services_test

import (
	"rally/services"
	"rally/stores"
	"rally/testutil"
)

// Can upvote up to a limit.
func (t *ServicesSuite) Test_VotingService_UpvotingUpperLimit() {
	u := t.MustCreateUser()
	b := t.MustCreateBoardWithVoteLimit(1)

	p := t.MustCreatePost(t.ValidPost(b, u))

	s := services.NewVotingService(stores.NewVotingStore(testutil.NewFakeStorage()), b.VotingStrategy)
	count, err := s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(1, count)

	postVotes, upvoted, err := s.Upvote(u, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(0, count)

	postVotes, upvoted, err = s.Upvote(u, p)
	t.NoError(err)
	t.False(upvoted)
	t.Equal(1, postVotes)
}

// Vote limit can be turned off.
func (t *ServicesSuite) Test_VotingService_UpvotingNoUpperLimit() {
	u := t.MustCreateUser()
	b := t.MustCreateBoardWithNoVoteLimit()
	p := t.MustCreatePost(t.ValidPost(b, u))

	s := services.NewVotingService(stores.NewVotingStore(testutil.NewFakeStorage()), b.VotingStrategy)

	_, err := s.VotesRemaining(u, b)
	t.Error(err) // No limit.

	postVotes, upvoted, err := s.Upvote(u, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	postVotes, upvoted, err = s.Upvote(u, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(2, postVotes)
}

// Boards keep independent vote counts.
func (t *ServicesSuite) Test_VotingService_BoardIsolation() {
	u := t.MustCreateUser()

	b1 := t.MustCreateBoardWithVoteLimit(1)
	b2 := t.MustCreateBoardWithVoteLimit(1)

	p1 := t.MustCreatePost(t.ValidPost(b1, u))
	p2 := t.MustCreatePost(t.ValidPost(b2, u))

	// It's fine to reuse the service, both boards have identical voting strategy.
	s := services.NewVotingService(stores.NewVotingStore(testutil.NewFakeStorage()), b1.VotingStrategy)

	postVotes, upvoted, err := s.Upvote(u, p1)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	count, err := s.VotesRemaining(u, b1)
	t.NoError(err)
	t.Equal(0, count)

	count, err = s.VotesRemaining(u, b2)
	t.NoError(err)
	t.Equal(1, count)

	postVotes, upvoted, err = s.Upvote(u, p2)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	count, err = s.VotesRemaining(u, b2)
	t.NoError(err)
	t.Equal(0, count)
}

// Can take votes back until going to zero.
func (t *ServicesSuite) Test_VotingService_Backsies() {
	u := t.MustCreateUser()
	b := t.MustCreateBoardWithVoteLimit(1)
	p := t.MustCreatePost(t.ValidPost(b, u))

	s := services.NewVotingService(stores.NewVotingStore(testutil.NewFakeStorage()), b.VotingStrategy)

	postVotes, upvoted, err := s.Upvote(u, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	count, err := s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(0, count)

	postVotes, downvoted, err := s.Downvote(u, p)
	t.NoError(err)
	t.True(downvoted)
	t.Equal(0, postVotes)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(1, count)

	postVotes, downvoted, err = s.Downvote(u, p)
	t.NoError(err)
	t.False(downvoted)
	t.Equal(0, postVotes)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(1, count)

	postVotes, upvoted, err = s.Upvote(u, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(0, count)

}

// Cannot take votes back for non-upvoted posts.
func (t *ServicesSuite) Test_VotingService_BacksiesOnlyForUpvotedPosts() {
	u := t.MustCreateUser()

	b := t.MustCreateBoardWithVoteLimit(1)

	p1 := t.MustCreatePost(t.ValidPost(b, u))
	p2 := t.MustCreatePost(t.ValidPost(b, u))

	s := services.NewVotingService(stores.NewVotingStore(testutil.NewFakeStorage()), b.VotingStrategy)

	postVotes, upvoted, err := s.Upvote(u, p1)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	count, err := s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(0, count)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(0, count)

	postVotes, downvoted, err := s.Downvote(u, p2)
	t.NoError(err)
	t.False(downvoted)
	t.Equal(0, postVotes)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(0, count)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(0, count)

	postVotes, downvoted, err = s.Downvote(u, p1)
	t.NoError(err)
	t.True(downvoted)
	t.Equal(0, postVotes)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(1, count)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(1, count)
}

// Votes can be refilled.
func (t *ServicesSuite) Test_VotingService_Refill() {
	u := t.MustCreateUser()
	b := t.MustCreateBoardWithVoteLimit(1)

	p := t.MustCreatePost(t.ValidPost(b, u))

	s := services.NewVotingService(stores.NewVotingStore(testutil.NewFakeStorage()), b.VotingStrategy)

	count, err := s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(1, count)

	postVotes, upvoted, err := s.Upvote(u, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(0, count)

	postVotes, upvoted, err = s.Upvote(u, p)
	t.NoError(err)
	t.False(upvoted)
	t.Equal(1, postVotes)

	err = s.Refill(b, *u)
	t.NoError(err)

	count, err = s.VotesRemaining(u, b)
	t.NoError(err)
	t.Equal(1, count)

	postVotes, upvoted, err = s.Upvote(u, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(2, postVotes)
}

// BUG: With two users up/downvoting a post, the post's votes were getting reset
// to the amount of votes each user placed INDIVIDUALLY.
func (t *ServicesSuite) Test_VotingService_Regression_MultipleUsers() {
	u1 := t.MustCreateUser()
	u2 := t.MustCreateUser()

	b := t.MustCreateBoardWithVoteLimit(10)

	p := t.MustCreatePost(t.ValidPost(b, u1))

	s := services.NewVotingService(stores.NewVotingStore(testutil.NewFakeStorage()), b.VotingStrategy)

	postVotes, upvoted, err := s.Upvote(u1, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	postVotes, upvoted, err = s.Upvote(u2, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(2, postVotes)

	count, err := s.VotesRemaining(u1, b)
	t.NoError(err)
	t.Equal(9, count)

	count, err = s.VotesRemaining(u2, b)
	t.NoError(err)
	t.Equal(9, count)

	postVotes, upvoted, err = s.Downvote(u2, p)
	t.NoError(err)
	t.True(upvoted)
	t.Equal(1, postVotes)

	count, err = s.VotesRemaining(u1, b)
	t.NoError(err)
	t.Equal(9, count)

	count, err = s.VotesRemaining(u2, b)
	t.NoError(err)
	t.Equal(10, count)
}
