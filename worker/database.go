package worker

import (
	"github.com/OHopiak/fractal-load-balancer/core"
	"github.com/jinzhu/gorm"
	//_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func (w *Worker) configureDatabase(dbConfig core.DatabaseConfig) {
	db, err := gorm.Open(dbConfig.Driver, dbConfig.ConnectionString)
	if err != nil {
		w.Server.Logger.Fatal("failed to connect database")
	}
	setupTables(db)
	w.db = db
}

func setupTables(db *gorm.DB) {
	core.SetupTables(db)
}
