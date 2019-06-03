package master

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	BaseData struct {
		LoggedIn bool
	}
)

func (m Master) homeHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "index.html", m.GetBaseData(c))
}

func (m Master) longProcessGetHandler(c echo.Context) error {
	return c.Render(http.StatusOK, "long_process.html", m.GetBaseData(c))
}

func (m Master) GetBaseData(c echo.Context) BaseData {
	userRaw := c.Get("user")
	return BaseData{
		LoggedIn: userRaw != nil,
	}
}