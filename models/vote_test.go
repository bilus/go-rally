package models_test

import "rally/models"

func (t *ModelSuite) Test_Vote() {
	p := t.ValidPost(t.MustCreateBoard(), t.MustCreateUser())
	verrs, err := t.DB.ValidateAndCreate(p)
	t.False(verrs.HasAny())
	t.NoError(err)

	v := models.Vote{
		PostID: p.ID,
		UserID: p.AuthorID,
	}
	verrs, err = t.DB.ValidateAndCreate(&v)
	t.False(verrs.HasAny())
	t.NoError(err)

	count, err := t.DB.Count("votes")
	t.NoError(err)
	t.Equal(1, count)

	t.DB.Destroy(p)

	count, err = t.DB.Count("votes")
	t.NoError(err)
	t.Equal(0, count)
}
