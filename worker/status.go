package worker

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	Status struct {
		Healthy bool   `json:"healthy"`
		Master  Master `json:"master"`
		Routes []*echo.Route `json:"routes"`
	}
)

func (w Worker) Status() *Status {
	return &Status{
		Healthy: true,
		Master:  *w.Master,
		Routes: w.Server.Routes(),
	}
}

func (w Worker) statusHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, w.Status())
}
