package core

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"net/url"
	"strconv"
	"strings"
)

type (
	Host struct {
		https bool
		IP    string `json:"ip"`
		Port  int    `json:"port"`
	}
)

func (h Host) String() string {
	result := h.IP
	if h.Port != 0 {
		result += fmt.Sprintf(":%d", h.Port)
	}
	return result
}

func (h Host) UrlRaw() string {
	result := ""
	if h.https {
		result += "https://"
	} else {
		result += "http://"
	}
	result += h.String()
	return result
}

func (h Host) Endpoint(path string) string {
	return h.UrlRaw() + path
}

func (h Host) Equal(other Host) bool {
	return h.IP == other.IP && h.Port == other.Port
}

func (h Host) Url() (*url.URL, error) {
	return url.Parse(h.UrlRaw())
}

func (h Host) ServerString() string {
	return fmt.Sprintf("%s:%d", h.IP, h.Port)
}

func HostFromString(host string) Host {
	data := strings.Split(host, ":")
	ip := ""
	port := 0
	if len(data) == 1 {
		ip = host
	} else {
		ip = data[0]
		lPort, _ := strconv.ParseInt(data[1], 10, 64)
		port = int(lPort)
	}
	return Host{
		IP: ip,
		Port: port,
	}
}

func StartServer(e *echo.Echo, host Host) {
	e.Logger.Fatal(e.Start(host.ServerString()))
}

func StartServerAsync(e *echo.Echo, host Host) chan error {
	err := make(chan error)
	go func() {
		err <- e.Start(host.ServerString())
	}()
	return err
}

func NewServer() *echo.Echo {
	e := echo.New()

	loggerFormat := `${time_rfc3339} ${remote_ip} [${method}] ${uri} ${status} ${error} ${latency_human}` + "\n"

	// Middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{Format: loggerFormat}))
	e.Use(middleware.Recover())

	// Configure Logger
	e.HideBanner = true
	e.Logger.SetLevel(log.INFO)
	e.Logger.SetHeader("${time_rfc3339} [${level}] ${long_file}:${line}")

	return e
}
