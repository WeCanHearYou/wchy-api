package middlewares

import (
	"github.com/getfider/fider/app"
	"github.com/getfider/fider/app/pkg/web"
)

// SetLocale defines given locale in context for all subsequent operations
func SetLocale(locale string) web.MiddlewareFunc {
	return func(next web.HandlerFunc) web.HandlerFunc {
		return func(c *web.Context) error {
			c.Set(app.LocaleCtxKey, locale)
			return next(c)
		}
	}
}
