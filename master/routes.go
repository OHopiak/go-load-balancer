package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (m *Master) routes() {
	m.balancer = NewWorkerBalancer(m.db)

	proxy := Proxy(m.balancer)
	userRequired := UserRequiredMiddleware()

	routes := core.RouteList{
		&core.BasicRoute{
			Path:       "/",
			Method:     core.GET,
			Controller: m.homeHandler,
		},
		&core.BasicRoute{
			Path:       "/status",
			Method:     core.GET,
			Controller: m.statusHandler,
		},
		&core.BasicRoute{
			Path:       "/worker/register",
			Method:     core.POST,
			Controller: m.registerWorkerHandler,
		},
		&core.BasicRoute{
			Path:       "/long_process",
			Method:     core.GET,
			Controller: m.longProcessGetHandler,
			Middleware: []echo.MiddlewareFunc{
				userRequired,
			},
		},
		&core.BasicRoute{
			Path:       "/long_process",
			Method:     core.POST,
			Controller: m.workerNotFound,
			Middleware: []echo.MiddlewareFunc{
				proxy,
			},
		},
		&core.BasicRoute{
			Path:       "/signup",
			Method:     core.GET,
			Controller: m.signUpGetHandler,
		},
		&core.BasicRoute{
			Path:       "/signup",
			Method:     core.POST,
			Controller: m.signUpPostHandler,
		},
		&core.BasicRoute{
			Path:       "/login",
			Method:     core.GET,
			Controller: m.loginGetHandler,
		},
		&core.BasicRoute{
			Path:       "/login",
			Method:     core.POST,
			Controller: m.loginPostHandler,
		},
		&core.BasicRoute{
			Path:       "/logout",
			Method:     core.GET,
			Controller: m.logoutGetHandler,
		},
		&core.BasicRoute{
			Path:       "/tasks",
			Method:     core.GET,
			Controller: m.taskViewHandler,
		},
		&core.GroupRoute{
			Path: "/task",
			Routes: core.RouteList{
				&core.BasicRoute{
					Path:       "/",
					Method:     core.GET,
					Controller: m.taskListHandler,
				},
				&core.BasicRoute{
					Path:       "/:id",
					Method:     core.GET,
					Controller: m.taskItemHandler,
				},
			},
		},
	}

	routes.Register(m.Server)
}

func (m *Master) middleware() {
	m.Server.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))
	m.Server.Use(UserMiddleware(m.db))
	m.Server.Static("/static", "master/static")
}

func (m *Master) templates() {

	t := core.NewTemplate("master/templates/").
		WithLayoutPath("master/templates/layout/").
		Parse(m.Server)
	m.Server.Renderer = t
}

func (m Master) workerNotFound(c echo.Context) error {
	return c.JSON(http.StatusNotFound, core.ErrorResponse{
		StatusCode: http.StatusNotFound,
		Message:    "no available workers found",
	})
}
