package actions

import (
	"net/http"
	"rally/models"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

func (t *ActionSuite) login(u *models.User) {
	t.Session.Set("current_user_id", u.ID)
}

func (t *ActionSuite) Authenticated(u *models.User) *models.User {
	t.login(u)
	return u
}

func (t *ActionSuite) stubbingCompleteUserAuth(u goth.User, err error, f func()) {
	oldHandler := gothic.CompleteUserAuth
	gothic.CompleteUserAuth = func(res http.ResponseWriter, req *http.Request) (goth.User, error) {
		return u, err
	}
	f()
	gothic.CompleteUserAuth = oldHandler
}

func (t *ActionSuite) Test_Auth_New() {
	res := t.HTML("/auth/new").Get()
	t.Equal(200, res.Code)
	t.Contains(res.Body.String(), "Sign In")
}

func (t *ActionSuite) Test_Auth_Create() {
	u := t.MustCreateUser()
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
		t.Run(tcase.Identifier, func() {
			res := t.HTML("/auth").Post(&models.User{
				Email:    tcase.Email,
				Password: tcase.Password,
			})

			t.Equal(tcase.Status, res.Code)
			t.Equal(tcase.RedirectURL, res.Location())
		})
	}
}

func (t *ActionSuite) Test_Auth_Redirect() {
	u := t.MustCreateUser()
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
		t.Run(tcase.identifier, func() {
			t.Session.Set("redirectURL", tcase.redirectURL)

			res := t.HTML("/auth").Post(u)

			t.Equal(302, res.Code)
			t.Equal(res.Location(), tcase.resultLocation)
		})
	}
}

func (t *ActionSuite) Test_Auth_Google() {
	res := t.HTML("/auth/google").Get()
	t.Equal(307, res.Code)
}

func (t *ActionSuite) Test_Auth_Callback_ExistingUser() {
	u := t.MustCreateUser()
	t.stubbingCompleteUserAuth(
		goth.User{
			UserID: u.GoogleUserID.String,
		},
		nil,
		func() {
			res := t.HTML("/auth/google/callback").Get()
			t.Equal(302, res.Code)
			t.Equal("/", res.Location())
		})
}

func (t *ActionSuite) Test_Auth_Callback_CreatesNewUser() {
	t.stubbingCompleteUserAuth(
		goth.User{
			Email:  "nosuchuser@example.com",
			UserID: "123",
		},
		nil,
		func() {
			res := t.HTML("/auth/google/callback").Get()
			t.Equal(302, res.Code)
			t.Equal("/", res.Location())
		})

	var u models.User
	err := t.DB.Where("email = ?", "nosuchuser@example.com").First(&u)
	t.NoError(err)
	t.Equal("nosuchuser@example.com", u.Email)
	t.True(u.GoogleUserID.Valid)
	t.Equal("123", u.GoogleUserID.String)
}
