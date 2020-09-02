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

func (ms *ModelSuite) SetupTest() {
	ms.Model.SetupTest()

	ms.LoadFixture("default")

	err := ms.DB.All(&ms.users)
	ms.NoError(err)

	ms.tempDir, err = ioutil.TempDir("", "rally")
	ms.NoError(err)
}

func (ms *ModelSuite) TearDownTest() {
	os.RemoveAll(ms.tempDir)
}

func Test_ModelSuite(t *testing.T) {
	model, err := suite.NewModelWithFixtures(packr.New("app:models:test:fixtures", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	ms := &ModelSuite{
		Model: model,
	}
	suite.Run(t, ms)
}
