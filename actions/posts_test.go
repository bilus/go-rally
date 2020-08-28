package actions

import (
	"rally/models"

	"github.com/Pallinder/go-randomdata"
)

func (as *ActionSuite) validPost(author *models.User) *models.Post {
	return &models.Post{
		Title:    randomdata.SillyName(),
		Body:     randomdata.Paragraph(),
		Votes:    randomdata.Number(1, 100),
		AuthorID: author.ID,
		Author:   *author,
	}
}

func (as *ActionSuite) emptPostDraft(author *models.User) *models.Post {
	p := as.validPost(author)
	p.Draft = true
	p.Title = ""
	p.Body = ""
	return p
}

// TODO: Create Fixtures struct, included in ModelSuite and ActionSuite.
func (as *ActionSuite) createPost(author *models.User, draft bool) *models.Post {
	p := as.validPost(author)
	p.Draft = draft
	verrs, err := as.DB.ValidateAndCreate(p)
	as.False(verrs.HasAny())
	as.NoError(err)
	return p
}

func (as *ActionSuite) createPosts(n int, author *models.User, drafts bool) []*models.Post {
	ps := make([]*models.Post, n)
	for i := range ps {
		ps[i] = as.createPost(author, drafts)
	}
	return ps
}

func (as *ActionSuite) Test_PostsResource_List_RequiresAuth() {
	res := as.HTML("/posts").Get()
	as.Equal(302, res.Code)
}

func (as *ActionSuite) Test_PostsResource_List() {
	originalCount, err := as.DB.Count("posts")
	as.NoError(err)

	u := as.authenticate()
	ps := as.createPosts(3, u, false)
	as.createPosts(1, u, true) // Drafts - hidden.
	res := as.HTML(as.PostsPath(nil)).Get()
	as.Equal(200, res.Code)

	doc := as.DOM(res)
	trs := doc.Find("tr.post-row")
	as.Equal(len(ps)+originalCount, trs.Length())
}

func (as *ActionSuite) Test_PostsResource_ListOwnDrafts() {
	originalCount, err := as.DB.Count("posts")
	as.NoError(err)

	u := as.authenticate()
	as.createPosts(3, u, false)          // Published - hidden
	drafts := as.createPosts(1, u, true) // Drafts - visible.
	u2, err := as.createUser()
	as.NoError(err)
	as.createPosts(2, u2, true) // Other user's drafts - hidden
	res := as.HTML(as.PostsPath(Opts{"drafts": true})).Get()
	as.Equal(200, res.Code)

	doc := as.DOM(res)
	trs := doc.Find("tr.post-row")
	as.Equal(len(drafts), trs.Length())
}

func (as *ActionSuite) Test_PostsResource_Show() {
	u := as.authenticate()
	p := as.createPost(u, false)
	res := as.HTML(as.PostPath(p)).Get()
	as.Contains(res.Body.String(), p.Title)
}

func (as *ActionSuite) Test_PostsResource_CreateDraft() {
	originalCount, err := as.DB.Count("posts")
	as.NoError(err)

	u := as.authenticate()
	p := as.emptPostDraft(u)

	res := as.HTML(as.PostsPath(nil)).Post(p)
	as.Equal(200, res.Code)

	count, err := as.DB.Count("posts")
	as.NoError(err)
	as.Equal(originalCount+1, count)

	var p1 models.Post
	as.NoError(as.DB.First(&p1))
	as.Equal(u.ID, p1.AuthorID)
}

func (as *ActionSuite) Test_PostsResource_Update() {
	u := as.authenticate()
	p := as.createPost(u, false)
	p.Title = "New title"
	res := as.HTML(as.PostPath(p)).Put(p)
	as.Equal(303, res.Code)
}

func (as *ActionSuite) Test_PostsResource_Update_OnlyAuthors() {
	as.authenticate()
	p := as.createPost(as.users[1], false) // Author != current user.
	p.Title = "New title"
	res := as.HTML(as.PostPath(p)).Put(p)
	as.Equal(401, res.Code)
}

func (as *ActionSuite) Test_PostsResource_Destroy() {
	originalCount, err := as.DB.Count("posts")
	as.NoError(err)

	u := as.authenticate()
	p := as.createPost(u, false)
	p.Title = "New title"
	res := as.HTML(as.PostPath(p)).Delete()
	as.Equal(303, res.Code)

	count, err := as.DB.Count("posts")
	as.NoError(err)
	as.Equal(originalCount, count)
}

func (as *ActionSuite) Test_PostsResource_Destroy_OnlyAuthors() {
	originalCount, err := as.DB.Count("posts")
	as.NoError(err)

	as.authenticate()
	p := as.createPost(as.users[1], false) // Author != current user.
	p.Title = "New title"
	res := as.HTML(as.PostPath(p)).Delete()
	as.Equal(401, res.Code)

	count, err := as.DB.Count("posts")
	as.NoError(err)
	as.Equal(originalCount+1, count)
}

func (as *ActionSuite) Test_PostsResource_Edit() {
	u := as.authenticate()
	p := as.createPost(u, false)
	res := as.HTML(as.EditPostPath(p)).Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_PostsResource_Edit_OnlyAuthors() {
	as.authenticate()
	p := as.createPost(as.users[1], false) // Author != current user.
	res := as.HTML(as.EditPostPath(p)).Get()
	as.Equal(401, res.Code)
}
