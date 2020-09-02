package actions

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

// func (as *ActionSuite) Test_BoardsResource_Create() {
// 	as.Fail("Not Implemented!")
// }

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
