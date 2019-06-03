package worker

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"net/http"
	"time"
)

type (
	LongProcessRequest struct {
		Iterations uint `json:"iterations"`
		SleepTimeMs uint `json:"sleep_time_ms"`
	}
	LongProcessResult struct {
		Status string `json:"status"`
	}
)

func (l LongProcessRequest) Validate() bool {
	return l.Iterations != 0 && l.SleepTimeMs != 0 && time.Duration(l.Iterations * l.SleepTimeMs) * time.Millisecond <= core.MaxTaskDuration
}

func (w Worker) theLongProcessItself(request LongProcessRequest, task core.Task) {
	if !request.Validate() {
		task.SetError(errors.New("the computing request is too big or invalid"))
		w.db.Save(&task)
		return
	}
	totalIterations := request.Iterations
	for i := uint(0); i < totalIterations; i++ {
		time.Sleep(time.Duration(request.SleepTimeMs) * time.Millisecond)
		task.DonePercentage = float32(i+1) / float32(totalIterations)
		w.db.Save(&task)
	}
	task.Result = []byte{1, 2, 3, 4, 5}
	result := LongProcessResult{
		Status: "success",
	}
	task.SaveResult(result)
	w.db.Save(&task)
}

func (w Worker) LongProcess(request LongProcessRequest) core.Task {
	task := w.NewTask()
	task.SaveParams(request)
	go w.theLongProcessItself(request, task)
	return task
}

func (w Worker) longProcessHandler(c echo.Context) error {
	request := new(LongProcessRequest)
	err := c.Bind(request)
	if err != nil {
		return err
	}

	task := w.LongProcess(*request)
	return c.JSON(http.StatusOK, task)
}
