package handlers_test

import (
	"testing"

	"net/http"
	"net/url"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/handlers"
	"github.com/getfider/fider/app/models"
	. "github.com/getfider/fider/app/pkg/assert"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/mock"
	"github.com/getfider/fider/app/pkg/oauth"
	"github.com/getfider/fider/app/pkg/web"
)

func TestSignOutHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/signout?redirect=/").
		AddCookie(web.CookieAuthName, "some-value").
		Execute(handlers.SignOut())

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("/")
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring(web.CookieAuthName + "=; Path=/; Expires=")
	Expect(response.Header().Get("Set-Cookie")).ContainsSubstring("Max-Age=0; HttpOnly")
}

func TestSignInByOAuthHandler(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.Execute(handlers.SignInByOAuth(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io/oauth/token?provider=facebook&redirect=")
}

func TestCallbackHandler_InvalidState(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, _ := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=abc").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusInternalServerError)
}

func TestCallbackHandler_InvalidCode(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()

	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://avengers.test.fider.io").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io")
}

func TestCallbackHandler_SignIn(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://avengers.test.fider.io&code=123").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io/oauth/facebook/token?code=123&path=")
}

func TestCallbackHandler_SignIn_WithPath(t *testing.T) {
	RegisterT(t)
	server, _ := mock.NewServer()

	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://avengers.test.fider.io/some-page&code=123").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://avengers.test.fider.io/oauth/facebook/token?code=123&path=%2Fsome-page")
}

func TestCallbackHandler_SignUp(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://login.test.fider.io/oauth/callback?state=http://demo.test.fider.io/signup&code=123").
		Execute(handlers.OAuthCallback(oauth.FacebookProvider))
	Expect(code).Equals(http.StatusTemporaryRedirect)

	location, _ := url.Parse(response.Header().Get("Location"))
	Expect(location.Host).Equals("demo.test.fider.io")
	Expect(location.Scheme).Equals("http")
	Expect(location.Path).Equals("/signup")
	ExpectOAuthToken(location.Query().Get("token"), &jwt.OAuthClaims{
		OAuthProvider: "facebook",
		OAuthID:       "FB123",
		OAuthName:     "Jon Snow",
		OAuthEmail:    "jon.snow@got.com",
	})
}

func TestOAuthTokenHandler_ExistingUserAndProvider(t *testing.T) {
	RegisterT(t)

	server, _ := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=123").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	ExpectFiderAuthCookie(response, mock.JonSnow)
}

func TestOAuthTokenHandler_NewUser(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=456&path=/hello").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io/hello")

	user, err := services.Users.GetByEmail("some.guy@facebook.com")
	Expect(err).IsNil()
	Expect(user.Name).Equals("Some Facebook Guy")

	ExpectFiderAuthCookie(response, user)
}

func TestOAuthTokenHandler_NewUserWithoutEmail(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.Users.Register(&models.User{
		Name:   "Some Guy",
		Email:  "",
		Tenant: mock.DemoTenant,
		Providers: []*models.UserProvider{
			&models.UserProvider{UID: "GO999", Name: oauth.GoogleProvider},
		},
	})

	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=798").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	user, err := services.Users.GetByID(3)
	Expect(err).IsNil()
	Expect(user.ID).Equals(3)
	Expect(user.Name).Equals("Some Guy")
	Expect(user.Providers).HasLen(1)

	user, err = services.Users.GetByID(4)
	Expect(err).IsNil()
	Expect(user.ID).Equals(4)
	Expect(user.Name).Equals("Mark")
	Expect(user.Providers).HasLen(1)

	Expect(code).Equals(http.StatusTemporaryRedirect)

	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	ExpectFiderAuthCookie(response, &models.User{
		ID:   4,
		Name: "Mark",
	})
}

func TestOAuthTokenHandler_ExistingUser_WithoutEmail(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	services.Users.Register(&models.User{
		Name:   "Some Facebook Guy",
		Email:  "",
		Tenant: mock.DemoTenant,
		Providers: []*models.UserProvider{
			&models.UserProvider{UID: "FB456", Name: oauth.FacebookProvider},
		},
	})

	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=456").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)

	_, err := services.Users.GetByID(4)
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)

	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	ExpectFiderAuthCookie(response, &models.User{
		ID:   3,
		Name: "Some Facebook Guy",
	})
}

func TestCallbackHandler_ExistingUser_NewProvider(t *testing.T) {
	RegisterT(t)

	server, services := mock.NewServer()
	code, response := server.
		WithURL("http://demo.test.fider.io/oauth/facebook/token?code=123").
		OnTenant(mock.DemoTenant).
		Execute(handlers.OAuthToken(oauth.GoogleProvider))

	Expect(code).Equals(http.StatusTemporaryRedirect)

	user, err := services.Users.GetByEmail("jon.snow@got.com")
	Expect(err).IsNil()
	Expect(user.Providers).HasLen(2)

	Expect(response.Header().Get("Location")).Equals("http://demo.test.fider.io")
	ExpectFiderAuthCookie(response, mock.JonSnow)
}

func TestCallbackHandler_NewUser_PrivateTenant(t *testing.T) {
	RegisterT(t)
	server, services := mock.NewServer()
	mock.AvengersTenant.IsPrivate = true

	code, response := server.
		WithURL("http://ideas.theavengers.com/oauth/facebook/token?code=456").
		OnTenant(mock.AvengersTenant).
		Execute(handlers.OAuthToken(oauth.FacebookProvider))

	user, err := services.Users.GetByEmail("some.guy@facebook.com")
	Expect(errors.Cause(err)).Equals(app.ErrNotFound)
	Expect(user).IsNil()

	Expect(code).Equals(http.StatusTemporaryRedirect)
	Expect(response.Header().Get("Location")).Equals("http://ideas.theavengers.com/not-invited")
	ExpectFiderAuthCookie(response, nil)
}

func ExpectOAuthToken(token string, expected *jwt.OAuthClaims) {
	user, err := jwt.DecodeOAuthClaims(token)
	Expect(err).IsNil()
	Expect(user.OAuthID).Equals(expected.OAuthID)
	Expect(user.OAuthName).Equals(expected.OAuthName)
	Expect(user.OAuthEmail).Equals(expected.OAuthEmail)
	Expect(user.OAuthProvider).Equals(expected.OAuthProvider)
}
