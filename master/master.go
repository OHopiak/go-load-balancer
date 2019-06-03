package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type (
	Master struct {
		core.Host
		Server   *echo.Echo
		balancer ProxyBalancer
		db       *gorm.DB
	}
)

func New(dbConfig core.DatabaseConfig) Master {
	m := Master{
		Host: core.Host{
			IP:   "localhost",
			Port: 8000,
		},
		Server: core.NewServer(),
	}
	m.configureDatabase(dbConfig)
	m.templates()
	m.middleware()
	m.routes()
	return m
}

func (m Master) WithIP(IP string) Master {
	m.IP = IP
	return m
}

func (m Master) WithPort(port int) Master {
	m.Port = port
	return m
}

func (m *Master) Start() {
	errChan := core.StartServerAsync(m.Server, m.Host)
	err := m.PostStart()
	if err != nil {
		m.Server.Logger.Fatal(err)
	}
	m.Server.Logger.Fatal(<-errChan)

}

func (m *Master) PostStart() error {
	for _, worker := range m.Workers() {
		m.Server.Logger.Infof("Adding worker %s", worker.Host())
		url, err := worker.Host().Url()
		if err != nil {
			return err
		}

		m.balancer.AddTarget(&ProxyTarget{
			Name: url.String(),
			URL:  url,
			Meta: echo.Map{
				"id": worker.ID,
			},
		})
	}
	return nil
}
