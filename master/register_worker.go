package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (m *Master) RegisterWorker(request *core.RegisterWorkerRequest, ip string) *core.RegisterWorkerResponse {
	worker, err := m.AddWorker(core.Host{
		IP:   ip,
		Port: request.Port,
	})

	if err != nil {
		return &core.RegisterWorkerResponse{
			Status: "FAILED",
		}
	}

	return &core.RegisterWorkerResponse{
		Status: "OK",
		WorkerId: worker.ID,
	}
}

func (m *Master) registerWorkerHandler(c echo.Context) error {
	request := new(core.RegisterWorkerRequest)
	err := c.Bind(request)
	if err != nil {
		return err
	}

	return c.JSON(http.StatusCreated, m.RegisterWorker(request, c.RealIP()))
}
