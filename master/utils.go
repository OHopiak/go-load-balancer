package master

import (
	"github.com/gorilla/sessions"
	"github.com/labstack/echo/v4"
)

func GetSession(c echo.Context) (*sessions.Session, error) {
	return GetSessionByName("session", c)
}

func GetSessionByName(name string, c echo.Context) (*sessions.Session, error) {
	store, ok := c.Get("_session_store").(sessions.Store)
	if !ok {
		c.Logger().Fatal("session is not configured")
	}
	return store.Get(c.Request(), name)
}