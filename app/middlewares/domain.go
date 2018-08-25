package middlewares

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/getfider/fider/app/pkg/validate"

	"github.com/getfider/fider/app/models"

	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/env"
	"github.com/getfider/fider/app/pkg/errors"
	"github.com/getfider/fider/app/pkg/jwt"
	"github.com/getfider/fider/app/pkg/web"
	"github.com/getfider/fider/app/pkg/web/util"
)

// Tenant adds either SingleTenant or MultiTenant to the pipeline
func Tenant() web.MiddlewareFunc {
	if env.IsSingleHostMode() {
		return SingleTenant()
	}
	return MultiTenant()
}

// SingleTenant inject default tenant into current context
func SingleTenant() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			tenant, err := c.Services().Tenants.First()
			if err != nil && errors.Cause(err) != app.ErrNotFound {
				return c.Failure(err)
			}

			if tenant != nil {
				c.SetTenant(tenant)
			}

			return next(c)
		}
	}
}

// MultiTenant extract tenant information from hostname and inject it into current context
func MultiTenant() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			hostname := c.Request.URL.Hostname()

			// If no tenant is specified, redirect user to getfider.com
			// This is only valid for fider.io hosting
			if (env.IsProduction() && hostname == "fider.io") ||
				(env.IsDevelopment() && hostname == "dev.fider.io") {
				if c.Request.URL.Path == "" || c.Request.URL.Path == "/" {
					return c.Redirect("https://getfider.com")
				}
			}

			tenant, err := c.Services().Tenants.GetByDomain(hostname)
			if err != nil && errors.Cause(err) != app.ErrNotFound {
				return c.Failure(err)
			}

			if tenant != nil {
				c.SetTenant(tenant)

				if tenant.CNAME != "" && !c.IsAjax() {
					baseURL := c.TenantBaseURL(tenant)
					if baseURL != c.BaseURL() {
						link := baseURL + c.Request.URL.RequestURI()
						c.Response.Header().Set("Link", fmt.Sprintf("<%s>; rel=\"canonical\"", link))
					}
				}
			}

			return next(c)
		}
	}
}

// RequireTenant returns 404 if tenant is not available
func RequireTenant() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			if c.Tenant() == nil {
				if env.IsSingleHostMode() {
					return c.Redirect("/signup")
				}
				return c.NotFound()
			}
			return next(c)
		}
	}
}

// User gets JWT Auth token from cookie and insert into context
func User() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			var (
				token string
				user  *models.User
			)

			cookie, err := c.Request.Cookie(web.CookieAuthName)
			if err == nil {
				token = cookie.Value
			} else {
				token = webutil.GetSignUpAuthCookie(c)
				if token != "" {
					webutil.AddAuthTokenCookie(c, token)
				}
			}

			if token != "" {
				claims, err := jwt.DecodeFiderClaims(token)
				if err != nil {
					c.RemoveCookie(web.CookieAuthName)
					c.SetClaims(claims)
					return next(c)
				}

				user, err = c.Services().Users.GetByID(claims.UserID)
				if err != nil {
					if errors.Cause(err) == app.ErrNotFound {
						c.RemoveCookie(web.CookieAuthName)
						return next(c)
					}
					return err
				}
			} else {
				authHeader := c.Request.GetHeader("Authorization")
				parts := strings.Split(authHeader, "Bearer")
				if len(parts) == 2 {
					apiKey := strings.TrimSpace(parts[1])
					user, err = c.Services().Users.GetByAPIKey(apiKey)
					if err != nil {
						if errors.Cause(err) == app.ErrNotFound {
							return c.HandleValidation(validate.Failed("API Key is invalid"))
						}
						return err
					}

					if !user.IsCollaborator() {
						return c.HandleValidation(validate.Failed("API Key is invalid"))
					}

					if impersonateUserIDStr := c.Request.GetHeader("X-Fider-UserID"); impersonateUserIDStr != "" {
						if !user.IsAdministrator() {
							return c.HandleValidation(validate.Failed("Only Administrators are allowed to impersonate another user"))
						}
						impersonateUserID, err := strconv.Atoi(impersonateUserIDStr)
						if err != nil {
							return c.HandleValidation(validate.Failed(fmt.Sprintf("User not found for given impersonate UserID '%s'", impersonateUserIDStr)))
						}
						user, err = c.Services().Users.GetByID(impersonateUserID)
						if err != nil {
							if errors.Cause(err) == app.ErrNotFound {
								return c.HandleValidation(validate.Failed(fmt.Sprintf("User not found for given impersonate UserID '%s'", impersonateUserIDStr)))
							}
							return err
						}
					}
				}
			}

			if user != nil && c.Tenant() != nil && user.Tenant.ID == c.Tenant().ID {
				c.SetUser(user)
			}

			return next(c)
		}
	}
}

// OnlyActiveTenants blocks requests for inactive tenants
func OnlyActiveTenants() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			if c.Tenant().Status == models.TenantActive {
				return next(c)
			}
			return c.NotFound()
		}
	}
}

// CheckTenantPrivacy blocks requests of unauthenticated users for private tenants
func CheckTenantPrivacy() web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c web.Context) error {
			if c.Tenant().IsPrivate && !c.IsAuthenticated() {
				return c.Redirect("/signin")
			}
			return next(c)
		}
	}
}
