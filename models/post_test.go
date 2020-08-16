package models

// NOTE: At least one user must exist.
func (ms *ModelSuite) validPost() *Post {
	var u User
	ms.NoError(ms.DB.First(&u))
	return &Post{
		Title:    "Lorem",
		Body:     "Ipsum.",
		AuthorID: u.ID,
	}
}

func (ms *ModelSuite) Test_Post_Create() {
	count, err := ms.DB.Count("posts")
	ms.NoError(err)
	ms.Equal(0, count)

	p := ms.validPost()

	verrs, err := ms.DB.ValidateAndCreate(p)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	count, err = ms.DB.Count("posts")
	ms.NoError(err)
	ms.Equal(1, count)
}

func (ms *ModelSuite) Test_Post_RequiresTitle() {
	p := ms.validPost()
	p.Title = ""

	verrs, err := ms.DB.ValidateAndCreate(p)
	ms.NoError(err)
	ms.True(verrs.HasAny())
}

func (ms *ModelSuite) Test_Post_RequiresBody() {
	p := ms.validPost()
	p.Body = ""

	verrs, err := ms.DB.ValidateAndCreate(p)
	ms.NoError(err)
	ms.True(verrs.HasAny())
}
