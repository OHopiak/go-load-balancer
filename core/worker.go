package core

import (
	"time"
)

type (
	Worker struct {
		ID        uint       `json:"id" gorm:"primary_key"`
		CreatedAt time.Time  `json:"created_at"`
		UpdatedAt time.Time  `json:"updated_at"`
		DeletedAt *time.Time `json:"deleted_at,omitempty" sql:"index"`

		IP   string `json:"ip"`
		Port int    `json:"port"`

		Healthy bool `json:"healthy"`
	}
)

func (w Worker) Host() Host {
	return Host{
		IP: w.IP,
		Port: w.Port,
	}
}