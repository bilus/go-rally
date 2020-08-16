package actions

import (
	"fmt"
	"rally/models"
)

func (as *ActionSuite) Test_Votes_Create() {
	u := as.authenticate()
	p := as.createPost(u)

	res := as.JavaScript(fmt.Sprintf("/posts/%v/votes", p.ID)).Post(nil)
	as.Equal(200, res.Code)

	var p1 models.Post
	err := as.DB.Find(&p1, p.ID)
	as.NoError(err)
	as.Equal(p.Votes+1, p1.Votes)

	u1 := models.User{}
	err = as.DB.Find(&u1, u.ID)
	as.NoError(err)
	as.Equal(u.Votes-1, u1.Votes)
}

func (as *ActionSuite) Test_Votes_Destroy() {
	u := as.authenticate()
	p := as.createPost(u)

	res := as.JavaScript(fmt.Sprintf("/posts/%v/votes", p.ID)).Post(nil)
	as.Equal(200, res.Code)

	res = as.JavaScript(fmt.Sprintf("/posts/%v/votes", p.ID)).Delete()
	as.Equal(200, res.Code)

	var p1 models.Post
	err := as.DB.Find(&p1, p.ID)
	as.NoError(err)
	as.Equal(p.Votes, p1.Votes)

	u1 := models.User{}
	err = as.DB.Find(&u1, u.ID)
	as.NoError(err)
	as.Equal(u.Votes, u1.Votes)
}

func (as *ActionSuite) Test_Votes_Destroy_OnlyUpvoted() {
	u := as.authenticate()
	p := as.createPost(u)

	res := as.JavaScript(fmt.Sprintf("/posts/%v/votes", p.ID)).Delete()
	as.Equal(422, res.Code)

	var p1 models.Post
	err := as.DB.Find(&p1, p.ID)
	as.NoError(err)
	as.Equal(p.Votes, p1.Votes)
}
