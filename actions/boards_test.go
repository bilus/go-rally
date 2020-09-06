package actions

import "rally/models"

// func (as *ActionSuite) Test_BoardsResource_List() {
// 	as.Fail("Not Implemented!")
// }

func (t *ActionSuite) Test_BoardsResource_Show_RequiresAuth() {
	b := t.MustCreateBoard()
	res := t.HTML(t.BoardPath(b)).Get()
	t.Equal(302, res.Code)
}

func (t *ActionSuite) Test_BoardResource_Show() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.MustCreateBoard()
	ps := t.MustCreatePosts(t.ValidPosts(1, b, u, false)...)
	res := t.HTML(t.BoardPath(b)).Get()
	t.Equal(200, res.Code)

	doc := t.DOM(res)
	trs := doc.Find("tr.post-row")
	t.Equal(len(ps), trs.Length())
}

func (t *ActionSuite) Test_BoardResource_ShowIncludesOwnDrafts() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.MustCreateBoard()
	published := t.MustCreatePosts(t.ValidPosts(3, b, u, false)...) // Published
	drafts := t.MustCreatePosts(t.ValidPosts(1, b, u, true)...)     // Own drafts
	u2 := t.MustCreateUser()
	t.MustCreatePosts(t.ValidPosts(5, b, u2, true)...) // Published by other user
	res := t.HTML(t.BoardPath(b)).Get()
	t.Equal(200, res.Code)

	doc := t.DOM(res)
	trs := doc.Find("tr.post-row")
	t.Equal(len(drafts)+len(published), trs.Length())
}

func (t *ActionSuite) Test_BoardsResource_Create() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.ValidBoardWithVoteLimit(1)
	res := t.HTML(t.BoardsPath()).Post(b)
	t.Equal(303, res.Code)

	// Makes creator the only member, and an owner at that.
	t.NoError(t.DB.First(b))
	member := &models.BoardMember{}
	t.NoError(t.DB.Where("board_id = ? AND user_id = ?", b.ID, u.ID).First(member))
	t.Equal(u.ID, member.UserID)
	t.True(member.IsOwner)
}

// func (as *ActionSuite) Test_BoardsResource_Update() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_BoardsResource_Destroy() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_BoardsResource_New() {
// 	as.Fail("Not Implemented!")
// }

// func (as *ActionSuite) Test_BoardsResource_Edit() {
// 	as.Fail("Not Implemented!")
// }
