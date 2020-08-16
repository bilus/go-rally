package models

import (
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

type ModelSuite struct {
	*suite.Model

	users []*User
}

func (as *ModelSuite) SetupTest() {
	as.Model.SetupTest()

	as.LoadFixture("default")

	err := as.DB.All(&as.users)
	as.NoError(err)
}

func Test_ModelSuite(t *testing.T) {
	model, err := suite.NewModelWithFixtures(packr.New("app:models:test:fixtures", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ModelSuite{
		Model: model,
	}
	suite.Run(t, as)
}
