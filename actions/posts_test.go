package actions

import (
	"rally/models"

	"github.com/Pallinder/go-randomdata"
	"github.com/PuerkitoBio/goquery"
)

func (as *ActionSuite) validPost() *models.Post {
	return &models.Post{
		Title: randomdata.SillyName(),
		Body:  randomdata.Paragraph(),
		Votes: randomdata.Number(1, 100),
	}
}

func (as *ActionSuite) createPost() *models.Post {
	p := as.validPost()
	verrs, err := as.DB.ValidateAndCreate(p)
	as.False(verrs.HasAny())
	as.NoError(err)
	return p
}

func (as *ActionSuite) createPosts(n int) []*models.Post {
	ps := make([]*models.Post, n)
	for i := range ps {
		ps[i] = as.createPost()
	}
	return ps
}

func (as *ActionSuite) Test_PostsResource_List_RequiresAuth() {
	res := as.HTML("/posts").Get()
	as.Equal(302, res.Code)
}

func (as *ActionSuite) Test_PostsResource_List() {
	as.authenticate()
	ps := as.createPosts(3)
	res := as.HTML(as.PostsPath()).Get()
	as.Equal(200, res.Code)
	doc := as.DOM(res)
	trs := doc.Find("tr.post-row")
	as.Equal(len(ps), trs.Length())
	trs.Each(func(i int, tr *goquery.Selection) {
		as.Equal(ps[i].Title, tr.Find("td.title").Text())
	})
}

func (as *ActionSuite) Test_PostsResource_Show() {
	as.authenticate()
	p := as.createPost()
	res := as.HTML(as.PostPath(p)).Get()
	as.Contains(res.Body.String(), p.Title)
}

func (as *ActionSuite) Test_PostsResource_Create() {
	as.authenticate()
	p := as.validPost()
	res := as.HTML(as.PostsPath()).Post(p)
	as.Equal(303, res.Code)

	count, err := as.DB.Count("posts")
	as.NoError(err)
	as.Equal(1, count)
}

func (as *ActionSuite) Test_PostsResource_Update() {
	as.authenticate()
	p := as.createPost()
	p.Title = "New title"
	res := as.HTML(as.PostPath(p)).Put(p)
	as.Equal(303, res.Code)
}

func (as *ActionSuite) Test_PostsResource_Destroy() {
	as.authenticate()
	p := as.createPost()
	p.Title = "New title"
	res := as.HTML(as.PostPath(p)).Delete()
	as.Equal(303, res.Code)

	count, err := as.DB.Count("posts")
	as.NoError(err)
	as.Equal(0, count)
}

func (as *ActionSuite) Test_PostsResource_New() {
	as.authenticate()
	res := as.HTML(as.NewPostsPath()).Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_PostsResource_Edit() {
	as.authenticate()
	p := as.createPost()
	res := as.HTML(as.EditPostPath(p)).Get()
	as.Equal(200, res.Code)
}
