package core

import (
	"github.com/labstack/echo/v4"
)

type (
	Method int

	Route interface {
		Register(router Router)
	}

	BasicRoute struct {
		Path       string
		Method     Method
		Controller echo.HandlerFunc
		Middleware []echo.MiddlewareFunc
	}

	GroupRoute struct {
		Path       string
		Middleware []echo.MiddlewareFunc
		Routes     RouteList
	}

	RouteList []Route

	Router interface {
		// CONNECT registers a new CONNECT route for a path with matching handler in the
		// router with optional route-level middleware.
		CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
		// DELETE registers a new DELETE route for a path with matching handler in the router
		// with optional route-level middleware.
		DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

		// GET registers a new GET route for a path with matching handler in the router
		// with optional route-level middleware.
		GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

		// HEAD registers a new HEAD route for a path with matching handler in the
		// router with optional route-level middleware.
		HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

		// OPTIONS registers a new OPTIONS route for a path with matching handler in the
		// router with optional route-level middleware.
		OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

		// PATCH registers a new PATCH route for a path with matching handler in the
		// router with optional route-level middleware.
		PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

		// POST registers a new POST route for a path with matching handler in the
		// router with optional route-level middleware.
		POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

		// PUT registers a new PUT route for a path with matching handler in the
		// router with optional route-level middleware.
		PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

		// TRACE registers a new TRACE route for a path with matching handler in the
		// router with optional route-level middleware.
		TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route

		// Group creates a new router group with prefix and optional group-level middleware.
		Group(prefix string, m ...echo.MiddlewareFunc) *echo.Group
	}
)

const (
	GET    Method = iota
	POST          = iota
	PUT           = iota
	PATCH         = iota
	DELETE        = iota
)

func (r *BasicRoute) Register(router Router) {
	switch r.Method {
	case GET:
		router.GET(r.Path, r.Controller, r.Middleware...)
	case POST:
		router.POST(r.Path, r.Controller, r.Middleware...)
	case DELETE:
		router.DELETE(r.Path, r.Controller, r.Middleware...)
	case PATCH:
		router.PATCH(r.Path, r.Controller, r.Middleware...)
	case PUT:
		router.PUT(r.Path, r.Controller, r.Middleware...)
	}
}

func (r *GroupRoute) Register(router Router) {
	group := router.Group(r.Path, r.Middleware...)
	r.Routes.Register(group)
}

func (r RouteList) Register(router Router) {
	for _, route := range r {
		route.Register(router)
	}
}