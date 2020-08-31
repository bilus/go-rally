package models

import (
	"io/ioutil"
	"strings"

	"github.com/gofrs/uuid"
)

func (ms *ModelSuite) Test_Attachment_Open() {
	expected := "Hello, world"
	a := Attachment{ID: uuid.Must(uuid.NewV4())}
	err := a.Save(strings.NewReader(expected))
	ms.NoError(err)

	r, err := a.Open()
	defer r.Close()
	actual, err := ioutil.ReadAll(r)
	ms.NoError(err)

	ms.Equal(expected, string(actual))
}
