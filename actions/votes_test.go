package actions

import (
	"fmt"
	"rally/models"
)

func (as *ActionSuite) createPost() *models.Post {
	p := &models.Post{
		Title: "Lorem",
		Body:  "Ipsum.",
	}

	verrs, err := as.DB.ValidateAndCreate(p)
	as.False(verrs.HasAny())
	as.NoError(err)
	return p
}

func (as *ActionSuite) Test_Votes_Create() {
	as.authenticate()
	p := as.createPost()

	res := as.JavaScript(fmt.Sprintf("/posts/%v/votes", p.ID)).Post(nil)
	as.Equal(200, res.Code)

	var p1 models.Post
	err := as.DB.Find(&p1, p.ID)
	as.NoError(err)
	as.Equal(1, p1.Votes)
}
