package actions

import (
	"rally/models"
)

func (t *ActionSuite) Test_PostsResource_Show() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.MustCreateBoard()
	p := t.MustCreatePost(t.ValidPost(b, u))
	res := t.HTML(t.BoardPostPath(p)).Get()
	t.Contains(res.Body.String(), p.Title)
}

func (t *ActionSuite) Test_PostsResource_CreateDraft() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.MustCreateBoard()

	res := t.HTML(t.BoardPostsPath(b)).Post(t.EmptyPostDraft(b, u))
	t.Equal(200, res.Code)

	count, err := t.DB.Count("posts")
	t.NoError(err)
	t.Equal(1, count)

	var p models.Post
	t.NoError(t.DB.First(&p))
	t.Equal(u.ID, p.AuthorID)
}

func (t *ActionSuite) Test_PostsResource_Update() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.MustCreateBoardWithVoteLimit(10)
	p := t.MustCreatePost(t.ValidPost(b, u))
	p.Title = "New title"
	res := t.HTML(t.BoardPostPath(p)).Put(p)
	t.Equal(303, res.Code)
}

func (t *ActionSuite) Test_PostsResource_Update_OnlyAuthors() {
	t.Authenticated(t.MustCreateUser())
	author := t.MustCreateUser()
	b := t.MustCreateBoard()
	p := t.MustCreatePost(t.ValidPost(b, author)) // Author != current user.
	p.Title = "New title"
	res := t.HTML(t.BoardPostPath(p)).Put(p)
	t.Equal(401, res.Code)
}

func (t *ActionSuite) Test_PostsResource_Destroy() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.MustCreateBoard()
	p := t.MustCreatePost(t.ValidPost(b, u))
	p.Title = "New title"
	res := t.HTML(t.BoardPostPath(p)).Delete()
	t.Equal(303, res.Code)

	count, err := t.DB.Count("posts")
	t.NoError(err)
	t.Equal(0, count)
}

func (t *ActionSuite) Test_PostsResource_Destroy_OnlyAuthors() {
	t.Authenticated(t.MustCreateUser())
	author := t.MustCreateUser()
	b := t.MustCreateBoard()
	p := t.MustCreatePost(t.ValidPost(b, author)) // Author != current user.
	p.Title = "New title"
	res := t.HTML(t.BoardPostPath(p)).Delete()
	t.Equal(401, res.Code)

	count, err := t.DB.Count("posts")
	t.NoError(err)
	t.Equal(1, count)
}

func (t *ActionSuite) Test_PostsResource_Edit() {
	current := t.Authenticated(t.MustCreateUser())
	t.MustCreateUser()
	b := t.MustCreateBoard()
	p := t.MustCreatePost(t.ValidPost(b, current)) // Author == current user.
	res := t.HTML(t.EditBoardPostPath(p)).Get()
	t.Equal(200, res.Code)
}

func (t *ActionSuite) Test_PostsResource_Edit_OnlyAuthors() {
	t.Authenticated(t.MustCreateUser())
	author := (t.MustCreateUser())
	b := t.MustCreateBoard()
	p := t.MustCreatePost(t.ValidPost(b, author)) // Author != current user.
	res := t.HTML(t.EditBoardPostPath(p)).Get()
	t.Equal(401, res.Code)
}
