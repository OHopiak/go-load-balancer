package core

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

)

type (
	DatabaseConfig struct {
		Driver string `json:"driver"`
		ConnectionString string `json:"connection_string"`
	}
)

func SetupTables(db *gorm.DB) {
	db.AutoMigrate(&Task{}, &Worker{}, &User{})
}
