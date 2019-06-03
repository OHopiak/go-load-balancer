package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/jinzhu/gorm"
	"github.com/labstack/echo/v4"
)

type (
	workerBalancer struct {
		*commonBalancer
		db *gorm.DB
		i  uint32
		busyWorkers uint
	}
)

// NewWorkerBalancer returns a random proxy balancer.
func NewWorkerBalancer(db *gorm.DB) ProxyBalancer {
	b := &workerBalancer{commonBalancer: newCommonBalancer(nil)}
	if db == nil {
		panic("db must be set before using this balancer")
	}
	b.db = db
	return b
}

// Next returns an upstream target using round-robin technique.
func (b *workerBalancer) Next(c echo.Context) *ProxyTarget {
	if len(b.targets) == 0 {
		return nil
	}
	for _, t := range b.targets {
		if b.workerCanAccept(t) {
			return t
		}
	}
	return nil
}

// AddTarget adds an upstream target to the list.
func (b *workerBalancer) AddTarget(target *ProxyTarget) bool {
	for _, t := range b.targets {
		if t.Name == target.Name {
			return false
		}
	}
	b.mutex.Lock()
	defer b.mutex.Unlock()
	b.targets = append(b.targets, target)
	return true
}

// RemoveTarget removes an upstream target from the list.
func (b *workerBalancer) RemoveTarget(name string) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()
	for i, t := range b.targets {
		if t.Name == name {
			b.targets = append(b.targets[:i], b.targets[i+1:]...)
			b.deleteWorker(t)
			return true
		}
	}
	return false
}

func (b *workerBalancer) deleteWorker(target *ProxyTarget) {
	rawId, ok := target.Meta["id"]
	if !ok {
		return
	}
	id := rawId.(uint)
	worker := core.Worker{}
	b.db.Delete(&worker, id)
}

func (b *workerBalancer) workerCanAccept(target *ProxyTarget) bool {
	rawId, ok := target.Meta["id"]
	if !ok {
		return false
	}
	id := rawId.(uint)
	var count uint
	b.db.Model(&core.Task{}).Where("worker_id=? and done_percentage < 1", id).Count(&count)
	return count < core.MaxTasks
}
