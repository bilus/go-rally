package models_test

import (
	"io/ioutil"
	"os"
	"rally/fixtures"
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

type ModelSuite struct {
	*suite.Model

	fixtures.Factory

	tempDir string
}

func (t *ModelSuite) SetupTest() {
	t.Model.SetupTest()
	t.Factory = fixtures.NewFactory(t.DB)

	var err error
	t.tempDir, err = ioutil.TempDir("", "rally")
	t.NoError(err)
}

func (t *ModelSuite) TearDownTest() {
	os.RemoveAll(t.tempDir)
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
