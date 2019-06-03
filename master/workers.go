package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/labstack/echo/v4"
)

func (m *Master) Workers() (workers []core.Worker) {
	m.db.Find(&workers)
	return
}

func (m *Master) AddWorker(host core.Host) (*core.Worker, error) {
	m.Server.Logger.Infof("Adding worker %s", host)
	for _, worker := range m.Workers() {
		if worker.Host().Equal(host) {
			worker.Healthy = true
			m.db.Save(worker)
			return nil, nil
		}
	}

	worker := core.Worker{
		Healthy: true,
		IP:    host.IP,
		Port:    host.Port,
	}

	url, err := worker.Host().Url()
	if err != nil {
		return nil, err
	}

	m.db.Create(&worker)
	m.balancer.AddTarget(&ProxyTarget{
		Name: url.String(),
		URL: url,
		Meta: echo.Map{
			"id": worker.ID,
		},
	})
	return &worker, nil
}

