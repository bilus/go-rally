package models

func validPost() *Post {
	return &Post{
		Title: "Lorem",
		Body:  "Ipsum.",
	}
}

func (ms *ModelSuite) Test_Post_Create() {
	count, err := ms.DB.Count("posts")
	ms.NoError(err)
	ms.Equal(0, count)

	p := validPost()

	verrs, err := ms.DB.ValidateAndCreate(p)
	ms.NoError(err)
	ms.False(verrs.HasAny())

	count, err = ms.DB.Count("posts")
	ms.NoError(err)
	ms.Equal(1, count)
}

func (ms *ModelSuite) Test_Post_RequiresTitle() {
	p := validPost()
	p.Title = ""

	verrs, err := ms.DB.ValidateAndCreate(p)
	ms.NoError(err)
	ms.True(verrs.HasAny())
}

func (ms *ModelSuite) Test_Post_RequiresBody() {
	p := validPost()
	p.Body = ""

	verrs, err := ms.DB.ValidateAndCreate(p)
	ms.NoError(err)
	ms.True(verrs.HasAny())
}
