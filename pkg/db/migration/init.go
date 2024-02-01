package migration

import (
	"github.com/hamiddarani/anakonda/internal/config"
	"github.com/hamiddarani/anakonda/pkg/db"
	"github.com/hamiddarani/anakonda/pkg/db/model"
	"github.com/hamiddarani/anakonda/pkg/logging"
	"gorm.io/gorm"
)

var logger = logging.NewLogger(config.Load(false).Logger)

func UP_1() {
	database := db.GetDb()

	createTables(database)
}

func createTables(database *gorm.DB) {
	tables := []interface{}{}

	tables = addNewTable(database, model.Task{}, tables)

	err := database.Migrator().CreateTable(tables...)
	if err != nil {
		logger.Error(logging.Postgres, logging.Migration, err.Error(), nil)
	}
	logger.Info(logging.Postgres, logging.Migration, "tables created", nil)
}

func addNewTable(database *gorm.DB, model interface{}, tables []interface{}) []interface{} {
	if !database.Migrator().HasTable(model) {
		tables = append(tables, model)
	}
	return tables
}
