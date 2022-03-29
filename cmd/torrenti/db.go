package main

import (
	"database/sql"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/pkg/errors"
	zlog "github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newDB(conf *DatabaseConf) (*sql.DB, *gorm.DB, error) {
	var err error
	var gdb *gorm.DB
	var db *sql.DB
	var dsn string
	var driver string

	switch conf.Type {
	case "sqlite":
		dbFile := conf.Database
		u := url.URL{
			Scheme: "file",
			Path:   dbFile,
		}
		dsn = u.String()
		driver = "sqlite"

		zlog.Debug().Str("dsn", dsn).Str("driver", driver).Msg("use sqlite")

		db, err = sql.Open(driver, dsn)
		if err != nil {
			return nil, nil, err
		}
		// sqlite
		db.SetMaxOpenConns(1)

		//if err = gdb.Exec("PRAGMA page_size = ?", 128*1024).Error; err != nil {
		//	return err
		//}
	default:
		err = errors.New("unsupported db type: " + conf.Type)
	}
	if err == nil {
		gdb, err = gorm.Open(sqlite.Dialector{
			Conn: db,
		}, &gorm.Config{
			Logger: newGLogger(conf),
		})
	}

	return db, gdb, err
}

func newGLogger(conf *DatabaseConf) logger.Interface {
	return logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), newGLoggerConf(conf))
}

func newGLoggerConf(conf *DatabaseConf) logger.Config {
	lc := logger.Config{
		SlowThreshold:             200 * time.Millisecond,
		LogLevel:                  logger.Warn,
		IgnoreRecordNotFoundError: conf.Log.IgnoreNotFound,
		Colorful:                  true,
	}
	if conf.Type == "sqlite" && conf.Log.SlowThreshold == 0 {
		lc.SlowThreshold = time.Second
	}
	return lc
}
