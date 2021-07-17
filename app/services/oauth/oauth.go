package oauth

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/models/cmd"
	"github.com/getfider/fider/app/models/dto"
	"github.com/getfider/fider/app/models/entity"
	"github.com/getfider/fider/app/models/enum"
	"github.com/getfider/fider/app/models/query"
	"github.com/getfider/fider/app/pkg/bus"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jsonq"
	"github.com/getfider/fider/app/pkg/validate"
	"github.com/getfider/fider/app/pkg/web"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
)

func init() {
	bus.Register(Service{})
}

type Service struct{}

func (s Service) Name() string {
	return "HTTP"
}

func (s Service) Category() string {
	return "OAuth"
}

func (s Service) Enabled() bool {
	return true
}

func (s Service) Init() {
	bus.AddHandler(parseOAuthRawProfile)
	bus.AddHandler(getOAuthAuthorizationURL)
	bus.AddHandler(getOAuthProfile)
	bus.AddHandler(getOAuthRawProfile)
	bus.AddHandler(listActiveOAuthProviders)
	bus.AddHandler(listAllOAuthProviders)
}

func getProviderStatus(key string) int {
	if key == "" {
		return enum.OAuthConfigDisabled
	}
	return enum.OAuthConfigEnabled
}

var (
	systemProviders = []*entity.OAuthConfig{
		{
			Provider:          app.FacebookProvider,
			DisplayName:       "Facebook",
			ProfileURL:        "https://graph.facebook.com/me?fields=name,email",
			Status:            getProviderStatus(env.Config.OAuth.Facebook.AppID),
			ClientID:          env.Config.OAuth.Facebook.AppID,
			ClientSecret:      env.Config.OAuth.Facebook.Secret,
			Scope:             "public_profile email",
			AuthorizeURL:      facebook.Endpoint.AuthURL,
			TokenURL:          facebook.Endpoint.TokenURL,
			JSONUserIDPath:    "id",
			JSONUserNamePath:  "name",
			JSONUserEmailPath: "email",
		},
		{
			Provider:          app.GoogleProvider,
			DisplayName:       "Google",
			ProfileURL:        "https://www.googleapis.com/oauth2/v2/userinfo",
			Status:            getProviderStatus(env.Config.OAuth.Google.ClientID),
			ClientID:          env.Config.OAuth.Google.ClientID,
			ClientSecret:      env.Config.OAuth.Google.Secret,
			Scope:             "profile email",
			AuthorizeURL:      "https://accounts.google.com/o/oauth2/v2/auth",
			TokenURL:          "https://www.googleapis.com/oauth2/v4/token",
			JSONUserIDPath:    "id",
			JSONUserNamePath:  "name",
			JSONUserEmailPath: "email",
		},
		{
			Provider:          app.GitHubProvider,
			DisplayName:       "GitHub",
			ProfileURL:        "https://api.github.com/user",
			Status:            getProviderStatus(env.Config.OAuth.GitHub.ClientID),
			ClientID:          env.Config.OAuth.GitHub.ClientID,
			ClientSecret:      env.Config.OAuth.GitHub.Secret,
			Scope:             "user:email",
			AuthorizeURL:      github.Endpoint.AuthURL,
			TokenURL:          github.Endpoint.TokenURL,
			JSONUserIDPath:    "id",
			JSONUserNamePath:  "name, login",
			JSONUserEmailPath: "email",
		},
	}
)

func parseOAuthRawProfile(ctx context.Context, c *cmd.ParseOAuthRawProfile) error {
	config, err := getConfig(ctx, c.Provider)
	if err != nil {
		return err
	}

	query := jsonq.New(c.Body)
	profile := &dto.OAuthUserProfile{
		ID:    strings.TrimSpace(query.String(config.JSONUserIDPath)),
		Name:  strings.TrimSpace(query.String(config.JSONUserNamePath)),
		Email: strings.ToLower(strings.TrimSpace(query.String(config.JSONUserEmailPath))),
	}

	if profile.ID == "" {
		return app.ErrUserIDRequired
	}

	if profile.Name == "" && profile.Email != "" {
		parts := strings.Split(profile.Email, "@")
		profile.Name = parts[0]
	}

	if profile.Name == "" {
		profile.Name = "Anonymous"
	}

	if len(validate.Email(ctx, profile.Email)) != 0 {
		profile.Email = ""
	}

	c.Result = profile
	return nil
}

func getOAuthAuthorizationURL(ctx context.Context, q *query.GetOAuthAuthorizationURL) error {
	config, err := getConfig(ctx, q.Provider)
	if err != nil {
		return err
	}

	oauthBaseURL := web.OAuthBaseURL(ctx)
	authURL, _ := url.Parse(config.AuthorizeURL)
	parameters := getProviderInitialParams(authURL)
	parameters.Add("client_id", config.ClientID)
	parameters.Add("scope", config.Scope)
	parameters.Add("redirect_uri", fmt.Sprintf("%s/oauth/%s/callback", oauthBaseURL, q.Provider))
	parameters.Add("response_type", "code")
	parameters.Add("state", q.Redirect+"|"+q.Identifier)

	authURL.RawQuery = parameters.Encode()
	q.Result = authURL.String()
	return nil
}

func getOAuthProfile(ctx context.Context, q *query.GetOAuthProfile) error {
	config, err := getConfig(ctx, q.Provider)
	if err != nil {
		return err
	}

	if config.Status == enum.OAuthConfigDisabled {
		return errors.New("Provider %s is disabled", q.Provider)
	}

	rawProfile := &query.GetOAuthRawProfile{Provider: q.Provider, Code: q.Code}
	err = bus.Dispatch(ctx, rawProfile)
	if err != nil {
		return err
	}

	parseRawProfile := &cmd.ParseOAuthRawProfile{Provider: q.Provider, Body: rawProfile.Result}
	err = bus.Dispatch(ctx, parseRawProfile)
	if err != nil {
		return err
	}

	q.Result = parseRawProfile.Result
	return nil
}

func getOAuthRawProfile(ctx context.Context, q *query.GetOAuthRawProfile) error {
	config, err := getConfig(ctx, q.Provider)
	if err != nil {
		return err
	}

	oauthBaseURL := web.OAuthBaseURL(ctx)
	exchange := (&oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.AuthorizeURL,
			TokenURL: config.TokenURL,
		},
		RedirectURL: fmt.Sprintf("%s/oauth/%s/callback", oauthBaseURL, q.Provider),
	}).Exchange

	oauthToken, err := exchange(ctx, q.Code)
	if err != nil {
		return err
	}

	if config.ProfileURL == "" {
		parts := strings.Split(oauthToken.AccessToken, ".")
		if len(parts) != 3 {
			return errors.New("AccessToken is not JWT")
		}

		body, _ := jwt.DecodeSegment(parts[1])
		q.Result = string(body)
		return nil
	}

	req := &cmd.HTTPRequest{
		URL:    config.ProfileURL,
		Method: "GET",
		Headers: map[string]string{
			"Authorization": "Bearer " + oauthToken.AccessToken,
		},
	}

	if err := bus.Dispatch(ctx, req); err != nil {
		return err
	}

	if req.ResponseStatusCode != 200 {
		return errors.New("Failed to request profile. Status Code: %d. Body: %s", req.ResponseStatusCode, string(req.ResponseBody))
	}

	q.Result = string(req.ResponseBody)
	return nil
}

func listActiveOAuthProviders(ctx context.Context, q *query.ListActiveOAuthProviders) error {
	allOAuthProviders := &query.ListAllOAuthProviders{}
	err := bus.Dispatch(ctx, allOAuthProviders)
	if err != nil {
		return err
	}

	list := make([]*dto.OAuthProviderOption, 0)
	for _, p := range allOAuthProviders.Result {
		if p.IsEnabled {
			list = append(list, p)
		}
	}
	q.Result = list
	return nil
}

func listAllOAuthProviders(ctx context.Context, q *query.ListAllOAuthProviders) error {
	oauthProviders := &query.ListCustomOAuthConfig{}
	err := bus.Dispatch(ctx, oauthProviders)
	if err != nil {
		return errors.Wrap(err, "failed to get list of custom OAuth providers")
	}

	oauthProviders.Result = append(oauthProviders.Result, systemProviders...)

	list := make([]*dto.OAuthProviderOption, 0)

	oauthBaseURL := web.OAuthBaseURL(ctx)
	for _, p := range oauthProviders.Result {
		list = append(list, &dto.OAuthProviderOption{
			Provider:         p.Provider,
			DisplayName:      p.DisplayName,
			ClientID:         p.ClientID,
			URL:              fmt.Sprintf("/oauth/%s", p.Provider),
			CallbackURL:      fmt.Sprintf("%s/oauth/%s/callback", oauthBaseURL, p.Provider),
			IsCustomProvider: string(p.Provider[0]) == "_",
			LogoBlobKey:      p.LogoBlobKey,
			IsEnabled:        p.Status == enum.OAuthConfigEnabled,
		})
	}

	q.Result = list
	return nil
}

func getConfig(ctx context.Context, provider string) (*entity.OAuthConfig, error) {
	for _, config := range systemProviders {
		if config.Status == enum.OAuthConfigEnabled && config.Provider == provider {
			return config, nil
		}
	}

	getCustomOAuth := &query.GetCustomOAuthConfigByProvider{Provider: provider}
	err := bus.Dispatch(ctx, getCustomOAuth)
	if err != nil {
		return nil, err
	}

	return getCustomOAuth.Result, nil
}
