package actions

import (
	"rally/models"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gobuffalo/httptest"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

type ActionSuite struct {
	*suite.Action
}

func (as ActionSuite) JavaScript(u string, args ...interface{}) *httptest.Request {
	r := httptest.New(as.App).HTML(u, args...)
	r.Headers["Accept"] = "text/javascript"
	return r
}

func (as ActionSuite) DOM(res *httptest.Response) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	as.NoError(err)
	return doc
}

func (as ActionSuite) Path(name string, opts map[string]interface{}) string {
	buildPath, ok := as.App.RouteHelpers()[name]
	as.True(ok)
	path, err := buildPath(opts)
	as.NoError(err)
	return string(path)
}

func (as ActionSuite) PostPath(p *models.Post) string {
	return as.Path("postPath", map[string]interface{}{"post_id": p.ID})
}

func (as ActionSuite) PostsPath() string {
	return as.Path("postsPath", nil)
}

func (as ActionSuite) NewPostsPath() string {
	return as.Path("newPostsPath", nil)
}

func (as ActionSuite) EditPostPath(p *models.Post) string {
	return as.Path("editPostPath", map[string]interface{}{"post_id": p.ID})
}

func Test_ActionSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(App(), packr.New("Test_ActionSuite", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ActionSuite{
		Action: action,
	}
	suite.Run(t, as)
}
