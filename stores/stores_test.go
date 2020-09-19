package stores

import (
	"rally/fixtures"
	"testing"

	"github.com/gobuffalo/packr/v2"
	"github.com/gobuffalo/suite/v3"
)

type StoresSuite struct {
	*suite.Model

	fixtures.Factory
}

func (t *StoresSuite) SetupTest() {
	t.Model.SetupTest()
	t.Factory = fixtures.NewFactory(t.DB)
}

func Test_StoresSuite(t *testing.T) {
	model, err := suite.NewModelWithFixtures(packr.New("app:models:test:fixtures", "../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	ms := &StoresSuite{
		Model: model,
	}
	suite.Run(t, ms)
}
