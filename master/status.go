package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

type (
	Status struct {
		Healthy bool          `json:"healthy"`
		Workers []core.Worker `json:"workers"`
	}
)

func (m *Master) Status() *Status{
	return &Status{
		Healthy: true,
		Workers: m.Workers(),
	}
}

func (m *Master) statusHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, m.Status())
}
