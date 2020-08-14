package actions

import (
	"testing"

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
