package models

func (ms *ModelSuite) Test_Vote() {
	initialCount, err := ms.DB.Count("votes")
	ms.NoError(err)

	p := ms.validPost()
	verrs, err := ms.DB.ValidateAndCreate(p)
	ms.False(verrs.HasAny())
	ms.NoError(err)

	v := Vote{
		PostID: p.ID,
		UserID: p.AuthorID,
	}
	verrs, err = ms.DB.ValidateAndCreate(&v)
	ms.False(verrs.HasAny())
	ms.NoError(err)

	count, err := ms.DB.Count("votes")
	ms.NoError(err)
	ms.Equal(initialCount+1, count)

	ms.DB.Destroy(p)

	count, err = ms.DB.Count("votes")
	ms.NoError(err)
	ms.Equal(initialCount, count)
}
