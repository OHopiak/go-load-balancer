package master

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
)

type (
	UserRequiredMiddlewareConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
	}
)

var (
	// DefaultUserRequiredMiddlewareConfig is the default user middleware config.
	DefaultUserRequiredMiddlewareConfig = UserRequiredMiddlewareConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

func UserRequiredMiddleware() echo.MiddlewareFunc {
	c := DefaultUserRequiredMiddlewareConfig
	return UserRequiredMiddlewareWithConfig(c)
}

func UserRequiredMiddlewareWithConfig(config UserRequiredMiddlewareConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultUserRequiredMiddlewareConfig.Skipper
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			user := c.Get("user")
			if user == nil {
				return c.Redirect(http.StatusMovedPermanently, "/login")
			}

			return next(c)
		}
	}
}
