package actions

import (
	"rally/fixtures"
	"rally/models"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/gobuffalo/httptest"
	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

type ActionSuite struct {
	*suite.Action
	fixtures.Factory
}

func (t *ActionSuite) SetupTest() {
	t.Action.SetupTest()

	t.Factory = fixtures.NewFactory(t.DB)
}

func (t ActionSuite) JavaScript(u string, args ...interface{}) *httptest.Request {
	r := httptest.New(t.App).HTML(u, args...)
	r.Headers["Accept"] = "text/javascript"
	return r
}

func (t ActionSuite) DOM(res *httptest.Response) *goquery.Document {
	doc, err := goquery.NewDocumentFromReader(res.Body)
	t.NoError(err)
	return doc
}

type Opts map[string]interface{}

func (t ActionSuite) Path(name string, opts Opts) string {
	buildPath, ok := t.App.RouteHelpers()[name]
	t.True(ok)
	path, err := buildPath(opts)
	t.NoError(err)
	return string(path)
}

func (t ActionSuite) BoardPostPath(p *models.Post) string {
	return t.Path("boardPostPath", Opts{"board_id": p.BoardID, "post_id": p.ID})
}

func (t ActionSuite) EditBoardPostPath(p *models.Post) string {
	return t.Path("editBoardPostPath", Opts{"board_id": p.BoardID, "post_id": p.ID})
}

func (t ActionSuite) BoardPostsPath(b *models.Board) string {
	return t.Path("boardPostsPath", Opts{"board_id": b.ID})
}

func (t ActionSuite) BoardPath(b *models.Board) string {
	return t.Path("boardPath", Opts{"board_id": b.ID})
}

func (t ActionSuite) BoardsPath() string {
	return t.Path("boardsPath", nil)
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
