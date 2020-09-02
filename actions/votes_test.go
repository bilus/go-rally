package actions

import (
	"fmt"
	"rally/models"
)

func (t *ActionSuite) Test_Votes_Create() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.MustCreateBoard()
	p := t.MustCreatePost(t.ValidPost(b, u))

	res := t.JavaScript(fmt.Sprintf("/posts/%v/votes", p.ID)).Post(nil)
	t.Equal(200, res.Code)

	var p1 models.Post
	err := t.DB.Find(&p1, p.ID)
	t.NoError(err)
	t.Equal(p.Votes+1, p1.Votes)

	u1 := models.User{}
	err = t.DB.Find(&u1, u.ID)
	t.NoError(err)
	t.Equal(u.Votes-1, u1.Votes)
}

func (t *ActionSuite) Test_Votes_Destroy() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.MustCreateBoard()
	p := t.MustCreatePost(t.ValidPost(b, u))

	res := t.JavaScript(fmt.Sprintf("/posts/%v/votes", p.ID)).Post(nil)
	t.Equal(200, res.Code)

	res = t.JavaScript(fmt.Sprintf("/posts/%v/votes", p.ID)).Delete()
	t.Equal(200, res.Code)

	var p1 models.Post
	err := t.DB.Find(&p1, p.ID)
	t.NoError(err)
	t.Equal(p.Votes, p1.Votes)

	u1 := models.User{}
	err = t.DB.Find(&u1, u.ID)
	t.NoError(err)
	t.Equal(u.Votes, u1.Votes)
}

func (t *ActionSuite) Test_Votes_Destroy_OnlyUpvoted() {
	u := t.Authenticated(t.MustCreateUser())
	b := t.MustCreateBoard()
	p := t.MustCreatePost(t.ValidPost(b, u))

	res := t.JavaScript(fmt.Sprintf("/posts/%v/votes", p.ID)).Delete()
	t.Equal(422, res.Code)

	var p1 models.Post
	err := t.DB.Find(&p1, p.ID)
	t.NoError(err)
	t.Equal(p.Votes, p1.Votes)
}
