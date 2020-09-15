package services_test

import (
	"rally/fixtures"
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

type ServicesSuite struct {
	*suite.Model

	fixtures.Factory
}

func (t *ServicesSuite) SetupTest() {
	t.Model.SetupTest()
	t.Factory = fixtures.NewFactory(t.DB)
}

func Test_ModelSuite(t *testing.T) {
	model, err := suite.NewModelWithFixtures(packr.New("app:models:test:fixtures", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	ms := &ServicesSuite{
		Model: model,
	}
	suite.Run(t, ms)
}
