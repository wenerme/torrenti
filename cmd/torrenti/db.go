package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"gorm.io/driver/postgres"

	"github.com/pkg/errors"
	zlog "github.com/rs/zerolog/log"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newDB(conf *DatabaseConf) (db *sql.DB, gdb *gorm.DB, err error) {
	dsn := conf.DSN
	driver := conf.Driver

	gc := &gorm.Config{
		Logger: newGLogger(conf),
	}
	switch conf.Type {
	case "sqlite":
		if driver == "" {
			driver = "sqlite"
		}
		if dsn == "" {
			dbFile := os.ExpandEnv(conf.Database)
			u := url.URL{
				Scheme: "file",
				Path:   dbFile,
			}
			for k, v := range conf.Attributes {
				u.Query().Add(k, v)
			}
			dsn = u.String()
		}

		zlog.Debug().Str("dsn", dsn).Str("driver", driver).Msg("use sqlite")

		db, err = sql.Open(driver, dsn)
		if err != nil {
			return
		}
		// sqlite
		db.SetMaxOpenConns(1)

		gdb, err = gorm.Open(sqlite.Dialector{
			Conn: db,
		}, gc)

		//if err = gdb.Exec("PRAGMA page_size = ?", 128*1024).Error; err != nil {
		//	return err
		//}

	case "postgres", "postgresql":
		if driver == "" {
			driver = "pgx"
		}

		if dsn == "" {
			o := make(map[string]string)
			o["user"] = conf.Username
			o["password"] = conf.Password
			o["host"] = conf.Host
			if conf.Port != 0 {
				o["port"] = fmt.Sprint(conf.Port)
			}
			o["dbname"] = conf.Database
			o["search_path"] = conf.Schema
			for k, v := range conf.Attributes {
				o[k] = v
			}
			sb := strings.Builder{}
			for k, v := range o {
				if v == "" {
					continue
				}
				sb.WriteString(k)
				sb.WriteRune('=')
				sb.WriteString(v)
				sb.WriteRune(' ')
			}
			dsn = sb.String()
			dsn = regexp.MustCompile(`\s+`).ReplaceAllString(dsn, " ")
			dsn = strings.TrimSpace(dsn)
		}

		db, err = sql.Open(driver, dsn)
		if err != nil {
			return
		}

		gdb, err = gorm.Open(postgres.New(postgres.Config{
			Conn: db,
		}), gc)
		if err != nil {
			return
		}

		if conf.Schema != "" && conf.CreateSchema {
			err = gdb.Exec("CREATE SCHEMA IF NOT EXISTS " + fmt.Sprintf("%q", conf.Schema)).Error
		}
	default:
		err = errors.New("unsupported db type: " + conf.Type)
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
	} else {
		lc.SlowThreshold = conf.Log.SlowThreshold
	}
	return lc
}
