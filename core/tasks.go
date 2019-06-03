package core

import (
	"encoding/json"
	"time"
)

type (
	TaskError string

	Task struct {
		ID             uint       `json:"id" gorm:"primary_key"`
		CreatedAt      time.Time  `json:"created_at"`
		UpdatedAt      time.Time  `json:"updated_at"`
		DeletedAt      *time.Time `json:"deleted_at,omitempty" sql:"index"`
		DonePercentage float32    `json:"done_percentage"`
		Params         []byte     `json:"params,omitempty"`
		Result         []byte     `json:"result,omitempty"`
		Error          *TaskError `json:"error,omitempty"`
		WorkerID       uint       `json:"-"`
		Worker         Worker     `json:"worker"`
	}

	ErrorResponse struct {
		StatusCode int    `json:"status_code"`
		Message    string `json:"message"`
	}
)

func (t *Task) SaveParams(params interface{}) {
	data, err := json.Marshal(params)
	if err != nil {
		t.SetError(err)
		return
	}
	t.Params = data
}

func (t Task) LoadParams(params interface{}) {
	err := json.Unmarshal(t.Params, params)
	if err != nil {
		t.SetError(err)
	}
}

func (t *Task) SaveResult(result interface{}) {
	data, err := json.Marshal(result)
	if err != nil {
		t.SetError(err)
		return
	}
	t.Result = data
}

func (t Task) LoadResult(result interface{}) {
	err := json.Unmarshal(t.Result, result)
	if err != nil {
		t.SetError(err)
	}
}

func (t *Task) SetError(err error) {
	t.Error = NewTaskError(err)
}

func (e TaskError) Error() string {
	return string(e)
}

func NewTaskError(err error) *TaskError {
	errorMsg := TaskError(err.Error())
	return &errorMsg
}
