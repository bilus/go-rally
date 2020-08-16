package actions

import (
	"net/http"
	"rally/models"

	"github.com/gobuffalo/nulls"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func (as *ActionSuite) createUser() (*models.User, error) {
	u := &models.User{
		Email:                "mark@example.com",
		Password:             "password",
		PasswordConfirmation: "password",
		GoogleUserID:         nulls.NewString("123"),
	}

	verrs, err := u.Create(as.DB)
	as.False(verrs.HasAny())

	return u, err
}

func (as *ActionSuite) login(u *models.User) {
	as.Session.Set("current_user_id", u.ID)
}

func (as *ActionSuite) authenticate() *models.User {
	u := as.users[0]
	as.login(u)
	return u
}

func (as *ActionSuite) stubbingCompleteUserAuth(u goth.User, err error, f func()) {
	oldHandler := gothic.CompleteUserAuth
	gothic.CompleteUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return u, err
	}
	f()
	gothic.CompleteUserAuth = oldHandler
}

func (as *ActionSuite) Test_Auth_New() {
	res := as.HTML("/auth/new").Get()
	as.Equal(200, res.Code)
	as.Contains(res.Body.String(), "Sign In")
}

func (as *ActionSuite) Test_Auth_Create() {
	u, err := as.createUser()
	as.NoError(err)

	tcases := []struct {
		Email       string
		Password    string
		Status      int
		RedirectURL string

		Identifier string
	}{
		{u.Email, u.Password, http.StatusFound, "/", "Valid"},
		{"noexist@example.com", "password", http.StatusUnauthorized, "", "Email Invalid"},
		{u.Email, "invalidPassword", http.StatusUnauthorized, "", "Password Invalid"},
	}

	for _, tcase := range tcases {
		as.Run(tcase.Identifier, func() {
			res := as.HTML("/auth").Post(&models.User{
				Email:    tcase.Email,
				Password: tcase.Password,
			})

			as.Equal(tcase.Status, res.Code)
			as.Equal(tcase.RedirectURL, res.Location())
		})
	}
}

func (as *ActionSuite) Test_Auth_Redirect() {
	u, err := as.createUser()
	as.NoError(err)

	tcases := []struct {
		redirectURL    interface{}
		resultLocation string

		identifier string
	}{
		{"/some/url", "/some/url", "RedirectURL defined"},
		{nil, "/", "RedirectURL nil"},
		{"", "/", "RedirectURL empty"},
	}

	for _, tcase := range tcases {
		as.Run(tcase.identifier, func() {
			as.Session.Set("redirectURL", tcase.redirectURL)

			res := as.HTML("/auth").Post(u)

			as.Equal(302, res.Code)
			as.Equal(res.Location(), tcase.resultLocation)
		})
	}
}

func (as *ActionSuite) Test_Auth_Google() {
	res := as.HTML("/auth/google").Get()
	as.Equal(307, res.Code)
}

func (as *ActionSuite) Test_Auth_Callback_ExistingUser() {
	u, err := as.createUser()
	as.NoError(err)
	as.stubbingCompleteUserAuth(
		goth.User{
			UserID: u.GoogleUserID.String,
		},
		nil,
		func() {
			res := as.HTML("/auth/google/callback").Get()
			as.Equal(302, res.Code)
			as.Equal("/", res.Location())
		})
}

func (as *ActionSuite) Test_Auth_Callback_CreatesNewUser() {
	as.stubbingCompleteUserAuth(
		goth.User{
			Email:  "nosuchuser@example.com",
			UserID: "123",
		},
		nil,
		func() {
			res := as.HTML("/auth/google/callback").Get()
			as.Equal(302, res.Code)
			as.Equal("/", res.Location())
		})

	var u models.User
	err := as.DB.Where("email = ?", "nosuchuser@example.com").First(&u)
	as.NoError(err)
	as.Equal("nosuchuser@example.com", u.Email)
	as.True(u.GoogleUserID.Valid)
	as.Equal("123", u.GoogleUserID.String)
}
