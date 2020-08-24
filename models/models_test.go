package models

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

type ModelSuite struct {
	*suite.Model

	users []*User

	tempDir string
}

func (as *ModelSuite) SetupTest() {
	as.Model.SetupTest()

	as.LoadFixture("default")

	err := as.DB.All(&as.users)
	as.NoError(err)

	as.tempDir, err = ioutil.TempDir("", "rally")
	as.NoError(err)
}

func (as *ModelSuite) TearDownTest() {
	os.RemoveAll(as.tempDir)
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
