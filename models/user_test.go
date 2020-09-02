package models_test

import "rally/models"

func (t *ModelSuite) Test_User_Create() {
	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	t.False(u.PasswordHash.Valid)

	verrs, err := u.Create(t.DB)
	t.NoError(err)
	t.False(verrs.HasAny())
	t.True(u.PasswordHash.Valid)
	t.NotZero(u.PasswordHash.String)

	count, err := t.DB.Count("users")
	t.NoError(err)
	t.Equal(1, count)
}

func (t *ModelSuite) Test_User_Create_ValidationErrors() {
	u := &models.User{
		Password: "password",
	}

	t.False(u.PasswordHash.Valid)

	verrs, err := u.Create(t.DB)
	t.NoError(err)
	t.True(verrs.HasAny())

	count, err := t.DB.Count("users")
	t.NoError(err)
	t.Equal(0, count)
}

func (t *ModelSuite) Test_User_Create_UserExists() {
	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
	}

	t.False(u.PasswordHash.Valid)

	verrs, err := u.Create(t.DB)
	t.NoError(err)
	t.False(verrs.HasAny())
	t.True(u.PasswordHash.Valid)

	count, err := t.DB.Count("users")
	t.NoError(err)
	t.Equal(1, count)

	u = &models.User{
		Email:    "mark@example.com",
		Password: "password",
	}

	verrs, err = u.Create(t.DB)
	t.NoError(err)
	t.True(verrs.HasAny())

	count, err = t.DB.Count("users")
	t.NoError(err)
	t.Equal(1, count)
}
