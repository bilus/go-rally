package actions

import (
	"rally/models"
)

func (as *ActionSuite) Test_Users_New() {
	res := as.HTML("/users/new").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_Users_Create() {
	initialCount, err := as.DB.Count("users")
	as.NoError(err)

	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
		Votes:                100,
	}

	res := as.HTML("/users").Post(u)
	as.Equal(302, res.Code)

	count, err := as.DB.Count("users")
	as.NoError(err)
	as.Equal(initialCount+1, count)
}
