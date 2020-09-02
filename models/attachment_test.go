package models_test

import (
	"io/ioutil"
	"rally/models"
	"strings"

	"github.com/gofrs/uuid"
)

func (t *ModelSuite) Test_Attachment_Open() {
	expected := "Hello, world"
	a := models.Attachment{ID: uuid.Must(uuid.NewV4())}
	err := a.Save(strings.NewReader(expected))
	t.NoError(err)

	r, err := a.Open()
	defer r.Close()
	actual, err := ioutil.ReadAll(r)
	t.NoError(err)

	t.Equal(expected, string(actual))
}
