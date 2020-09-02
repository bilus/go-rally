package models_test

func (t *ModelSuite) Test_Post_Create() {
	p := t.ValidPost(t.MustCreateBoard(), t.MustCreateUser())

	verrs, err := t.DB.ValidateAndCreate(p)
	t.NoError(err)
	t.False(verrs.HasAny())

	count, err := t.DB.Count("posts")
	t.NoError(err)
	t.Equal(1, count)
}

func (t *ModelSuite) Test_Post_RequiresTitle() {
	p := t.ValidPost(t.MustCreateBoard(), t.MustCreateUser())
	p.Title = ""

	verrs, err := t.DB.ValidateAndCreate(p)
	t.NoError(err)
	t.True(verrs.HasAny())
}

func (t *ModelSuite) Test_Post_RequiresBody() {
	p := t.ValidPost(t.MustCreateBoard(), t.MustCreateUser())
	p.Body = ""

	verrs, err := t.DB.ValidateAndCreate(p)
	t.NoError(err)
	t.True(verrs.HasAny())
}

func (t *ModelSuite) Test_Post_DraftRequiresNoTitleNorBody() {
	p := t.ValidPost(t.MustCreateBoard(), t.MustCreateUser())
	p.Draft = true
	p.Title = ""
	p.Body = ""

	verrs, err := t.DB.ValidateAndCreate(p)
	t.NoError(err)
	t.False(verrs.HasAny())
}
