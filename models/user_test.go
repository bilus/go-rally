package models

func (ms *ModelSuite) Test_User_Create() {
	initialCount, err := ms.DB.Count("users")
	ms.NoError(err)

	u := &User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	ms.False(u.PasswordHash.Valid)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.True(u.PasswordHash.Valid)
	ms.NotZero(u.PasswordHash.String)

	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(initialCount+1, count)
}

func (ms *ModelSuite) Test_User_Create_ValidationErrors() {
	initialCount, err := ms.DB.Count("users")
	ms.NoError(err)

	u := &User{
		Password: "password",
	}

	ms.False(u.PasswordHash.Valid)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())

	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(initialCount, count)
}

func (ms *ModelSuite) Test_User_Create_UserExists() {
	initialCount, err := ms.DB.Count("users")
	ms.NoError(err)

	u := &User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	ms.False(u.PasswordHash.Valid)

	verrs, err := u.Create(ms.DB)
	ms.NoError(err)
	ms.False(verrs.HasAny())
	ms.True(u.PasswordHash.Valid)

	count, err := ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(initialCount+1, count)

	u = &User{
		Email:    "mark@example.com",
		Password: "password",
	}

	verrs, err = u.Create(ms.DB)
	ms.NoError(err)
	ms.True(verrs.HasAny())

	count, err = ms.DB.Count("users")
	ms.NoError(err)
	ms.Equal(initialCount+1, count)
}
