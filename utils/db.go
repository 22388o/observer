package utils

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)
import "gorm.io/gorm"

var Gdb *gorm.DB

func InitDb(connstr string) {
	var err error
	newLogger := logger.New(
		log.New(os.Stderr, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,         // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		},
	)

	//db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	Gdb, err = gorm.Open(mysql.New(mysql.Config{
		DSN: connstr,
	}), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}
	// Migrate the schema
	Gdb.AutoMigrate()

	Gdb = Gdb.Debug()
}
