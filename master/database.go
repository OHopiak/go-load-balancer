package master

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/jinzhu/gorm"
)

func (m *Master) configureDatabase(dbConfig core.DatabaseConfig) {
	db, err := gorm.Open(dbConfig.Driver, dbConfig.ConnectionString)
	if err != nil {
		m.Server.Logger.Fatal("failed to connect database ", err)
	}
	setupTables(db)
	m.db = db
}

func setupTables(db *gorm.DB) {
	core.SetupTables(db)
}
