package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type (
	UserMiddlewareConfig struct {
		// Skipper defines a function to skip middleware.
		Skipper middleware.Skipper
		DB      *gorm.DB
	}
)

var (
	// DefaultUserMiddlewareConfig is the default user middleware config.
	DefaultUserMiddlewareConfig = UserMiddlewareConfig{
		Skipper: middleware.DefaultSkipper,
	}
)

func UserMiddleware(db *gorm.DB) echo.MiddlewareFunc {
	c := DefaultUserMiddlewareConfig
	c.DB = db
	return UserMiddlewareWithConfig(c)
}

func UserMiddlewareWithConfig(config UserMiddlewareConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = DefaultUserMiddlewareConfig.Skipper
	}

	if config.DB == nil {
		panic("db in user middleware must be set")
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if config.Skipper(c) {
				return next(c)
			}

			sess, err := GetSession(c)
			if err != nil {
				return err
			}

			userIdRaw, ok := sess.Values["user_id"]
			if !ok {
				return next(c)
			}

			userId, ok := userIdRaw.(uint)
			if !ok {
				c.Logger().Error("the user ID is not uint")
				return next(c)
			}

			if userId == 0 {
				c.Logger().Error("the user ID is set to 0")
				return next(c)
			}

			user := core.User{}
			config.DB.First(&user, userId)
			if user.ID == 0 {
				return next(c)
			}

			c.Set("user", user)

			return next(c)
		}
	}
}
