package stores

import (
	"rally/models"
	"rally/testutil"

	"github.com/johncgriffin/yogofn"
)

func (t StoresSuite) Test_ReactionStore_Empty() {
	u := t.MustCreateUser()
	b := t.MustCreateBoardWithVoteLimit(1)
	p := t.MustCreatePost(t.ValidPost(b, u))

	s := NewReactionStore(testutil.NewFakeStorage())
	reactions, err := s.ListReactionsToPost(p)
	t.NoError(err)
	t.Empty(reactions)
}

func (t StoresSuite) Test_ReactionStore_AddingReactions() {
	u := t.MustCreateUser()
	b := t.MustCreateBoardWithVoteLimit(1)
	p := t.MustCreatePost(t.ValidPost(b, u))

	s := NewReactionStore(testutil.NewFakeStorage())

	t.NoError(s.AddReactionToPost(u, p, "smile"))
	t.NoError(s.AddReactionToPost(u, p, "sad"))

	reactions, err := s.ListReactionsToPost(p)
	t.NoError(err)
	t.Len(reactions, 2)

	t.ElementsMatch([]string{"smile", "sad"},
		yogofn.Map(emoji, reactions).([]string))

	t.Len(reactions[0].UserIDs, 1)
	t.Len(reactions[1].UserIDs, 1)
}

func (t StoresSuite) Test_ReactionStore_AddingReactions_TwoUsers() {
	u1 := t.MustCreateUser()
	u2 := t.MustCreateUser()
	b := t.MustCreateBoardWithVoteLimit(99)
	p := t.MustCreatePost(t.ValidPost(b, u1))

	s := NewReactionStore(testutil.NewFakeStorage())

	t.NoError(s.AddReactionToPost(u1, p, "smile"))
	t.NoError(s.AddReactionToPost(u2, p, "smile"))
	t.NoError(s.AddReactionToPost(u1, p, "sad"))
	t.NoError(s.AddReactionToPost(u2, p, "happy"))

	reactions, err := s.ListReactionsToPost(p)
	t.NoError(err)
	t.Len(reactions, 3)

	t.ElementsMatch([]string{"smile", "sad", "happy"},
		yogofn.Map(emoji, reactions).([]string))

	t.True(yogofn.Any(func(r models.Reaction) bool {
		return r.Emoji == "smile" && len(r.UserIDs) == 2
	}, reactions))

	t.True(yogofn.Any(func(r models.Reaction) bool {
		return r.Emoji == "sad" && len(r.UserIDs) == 1
	}, reactions))

	t.True(yogofn.Any(func(r models.Reaction) bool {
		return r.Emoji == "happy" && len(r.UserIDs) == 1
	}, reactions))
}

func (t StoresSuite) Test_ReactionStore_AddingReactions_MaxOneReactionPerUser() {
	u := t.MustCreateUser()
	b := t.MustCreateBoardWithVoteLimit(1)
	p := t.MustCreatePost(t.ValidPost(b, u))

	s := NewReactionStore(testutil.NewFakeStorage())

	t.NoError(s.AddReactionToPost(u, p, "smile"))
	t.NoError(s.AddReactionToPost(u, p, "smile"))
	t.NoError(s.AddReactionToPost(u, p, "sad"))
	t.NoError(s.AddReactionToPost(u, p, "sad"))
	t.NoError(s.AddReactionToPost(u, p, "sad"))

	reactions, err := s.ListReactionsToPost(p)
	t.NoError(err)
	t.Len(reactions, 2)

	t.ElementsMatch([]string{"smile", "sad"},
		yogofn.Map(emoji, reactions).([]string))

	t.Len(reactions[0].UserIDs, 1)
	t.Len(reactions[1].UserIDs, 1)
}

func (t StoresSuite) Test_ReactionStore_RemovingReactions() {
	u1 := t.MustCreateUser()
	u2 := t.MustCreateUser()
	b := t.MustCreateBoardWithVoteLimit(99)
	p := t.MustCreatePost(t.ValidPost(b, u1))

	s := NewReactionStore(testutil.NewFakeStorage())

	t.NoError(s.AddReactionToPost(u1, p, "smile"))
	t.NoError(s.AddReactionToPost(u2, p, "smile"))
	t.NoError(s.AddReactionToPost(u1, p, "sad"))
	t.NoError(s.AddReactionToPost(u2, p, "happy"))

	t.NoError(s.RemoveReactionToPost(u1, p, "XXX"))
	t.NoError(s.RemoveReactionToPost(u1, p, "smile"))
	t.NoError(s.RemoveReactionToPost(u2, p, "happy"))

	reactions, err := s.ListReactionsToPost(p)
	t.NoError(err)
	t.Len(reactions, 2)

	t.Len(reactions[0].UserIDs, 1)
	t.Len(reactions[1].UserIDs, 1)
}

func emoji(r models.Reaction) string {
	return r.Emoji
}
