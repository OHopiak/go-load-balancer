package worker

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net"
	"time"
)

type (
	Master struct {
		core.Host `json:"host"`
		Connected bool `json:"connected"`
	}

	Worker struct {
		core.Host
		ID     *uint
		Master *Master
		Server *echo.Echo
		db     *gorm.DB
	}
)

func New(master core.Host, dbConfig core.DatabaseConfig) Worker {
	w := Worker{
		Host: core.Host{
			IP:   "localhost",
			Port: 0,
		},
		Server: core.NewServer(),
		Master: &Master{
			Host: master,
		},
		ID: new(uint),
	}

	w.configureDatabase(dbConfig)
	w.routes()
	return w
}

func (w Worker) WithIP(IP string) Worker {
	w.IP = IP
	return w
}

func (w Worker) WithPort(port int) Worker {
	w.Port = port
	return w
}

func (w *Worker) Start() {
	errChan := core.StartServerAsync(w.Server, w.Host)
	err := w.PostStart()
	if err != nil {
		w.Server.Logger.Fatal(err)
	}
	w.Server.Logger.Fatal(<-errChan)
}

func (w *Worker) PostStart() error {
	counter := 20
	for w.Server.Listener == nil && counter > 0 {
		time.Sleep(200 * time.Millisecond)
		counter--
	}
	if w.Server.Listener == nil {
		return errors.New("server start failed")
	}
	w.Port = w.Server.Listener.Addr().(*net.TCPAddr).Port

	err := w.Register()
	if err != nil {
		return err
	}
	return nil
}

func (w *Worker) NewTask() core.Task {
	worker := core.Worker{}
	w.db.First(&worker, *w.ID)
	task := core.Task{
		Worker: worker,
	}
	w.db.Create(&task)
	return task
}
