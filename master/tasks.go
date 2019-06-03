package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
)

type (
	TaskViewData struct {
		BaseData
		Tasks []core.Task
	}
)

func (m Master) taskListHandler(c echo.Context) error {
	var tasks []core.Task
	m.db.Find(&tasks)
	return c.JSON(http.StatusOK, tasks)
}

func (m Master) taskViewHandler(c echo.Context) error {
	var tasks []core.Task
	m.db.Preload("Worker").Find(&tasks)
	data := TaskViewData{
		BaseData: m.GetBaseData(c),
		Tasks: tasks,
	}
	return c.Render(http.StatusOK, "tasks.html", data)
}

func (m Master) taskItemHandler(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		return err
	}
	var task core.Task
	m.db.First(&task, id)
	if task.ID == 0 {
		return c.JSON(http.StatusNotFound, &core.ErrorResponse{
			StatusCode: http.StatusNotFound,
			Message: "task with the such ID is not found",
		})
	}
	return c.JSON(http.StatusOK, task)
}
