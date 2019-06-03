package worker

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/labstack/echo/v4/middleware"
)

func (w *Worker) routes() {
	routes := core.RouteList{
		&core.BasicRoute{
			Path:       "/status",
			Method:     core.GET,
			Controller: w.statusHandler,
		},
		&core.BasicRoute{
			Path:       "/long_process",
			Method:     core.POST,
			Controller: w.longProcessHandler,
		},
	}

	routes.Register(w.Server)
}

func (w *Worker) middleware() {
	w.Server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{w.Master.UrlRaw()},
	}))
}
