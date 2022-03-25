package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/wenerme/torrenti/pkg/indexer/models"
	"go.uber.org/multierr"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v3"

	"github.com/adrg/xdg"
	_ "github.com/glebarez/go-sqlite"
	"github.com/wenerme/torrenti/pkg/indexer"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	app := &cli.App{
		Name:  "torrenti",
		Usage: "torrent indexer",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Usage:   "config file",
				Value:   "",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Value: "",
			},
		},
		Before: setup,
		Commands: cli.Commands{
			{
				Name: "config",
				Subcommands: cli.Commands{
					{
						Name: "show",
						Action: func(ctx *cli.Context) error {
							return printYaml(_conf)
						},
					},
				},
			},
			{
				Name:   "stat",
				Action: showStat,
			},
			{
				Name: "torrent",
				Subcommands: cli.Commands{
					{
						Name:   "add",
						Usage:  "add to index",
						Action: addTorrent,
					},
				},
			},
			{
				Name: "magnet",
				Subcommands: cli.Commands{
					{
						Name:  "get-url",
						Usage: "get manget url of torrent",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func printYaml(v interface{}) error {
	bytes, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	fmt.Println(strings.TrimSpace(string(bytes)))
	return nil
}

var (
	_dataDir = filepath.Join(xdg.DataHome, "torrenti")
	_conf    = &Config{
		DataDir: _dataDir,
	}
)

func setup(ctx *cli.Context) error {
	conf := _conf

	{
		l := zerolog.InfoLevel
		if lvl, err := zerolog.ParseLevel(ctx.String("log-level")); err == nil {
			l = lvl
		}
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		log.Logger = zerolog.New(output).Level(l).With().Timestamp().Logger()
	}

	cfgFile := os.ExpandEnv(ctx.String("config"))
	if cfgFile == "" {
		cfgFile = filepath.Join(xdg.ConfigHome, "torrenti", "config.yaml")
		if _, err := os.Stat(cfgFile); err != nil {
			log.Debug().Str("path", cfgFile).Msg("default config file not found")
			cfgFile = ""
		}
	}

	if cfgFile != "" {
		log.Debug().Str("path", cfgFile).Msg("loading config file")

		bytes, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			return err
		}
		err = yaml.Unmarshal(bytes, conf)
		if err != nil {
			return err
		}
	}

	if err := env.Parse(conf); err != nil {
		return err
	}

	if ctx.String("log-level") != "" {
		conf.Log.Level = ctx.String("log-level")
	}
	if l, err := zerolog.ParseLevel(conf.Log.Level); err == nil {
		log.Logger = log.Level(l)
	}

	_dataDir = conf.DataDir
	log.Debug().Str("data_dir", _dataDir).Msg("conf")
	log.Debug().Str("log_level", conf.Log.Level).Msg("conf")

	err := os.MkdirAll(conf.DataDir, 0o755)
	if err != nil {
		return err
	}
	_ctx = ctx

	var gdb *gorm.DB
	var db *sql.DB
	var dsn string
	var driver string
	switch conf.DB.Type {
	case "sqlite":
		dbFile := conf.DB.Database
		if conf.DB.Database == "" {
			dbFile = filepath.Join(_dataDir, "torrenti.sqlite")
		}
		u := url.URL{
			Scheme: "file",
			Path:   dbFile,
		}
		dsn = u.String()
		driver = "sqlite"

		log.Debug().Str("dsn", dsn).Str("driver", driver).Msg("use sqlite")

		db, err = sql.Open("sqlite", dsn)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
		gdb, err = gorm.Open(sqlite.Dialector{
			Conn: db,
		})
	default:
		err = errors.New("unsupported db type: " + conf.DB.Type)
	}
	if err != nil {
		return err
	}

	_indexer, err = indexer.NewIndexer(indexer.NewIndexerOptions{DB: gdb})
	return err
}

func addTorrent(ctx *cli.Context) error {
	idx := getIndexer()
	for _, v := range ctx.Args().Slice() {
		log.Info().Str("torrent", v).Msg("add torrent")
		t, err := indexer.ParseTorrent(v)
		if err != nil {
			return err
		}
		err = t.Load()
		if err != nil {
			return err
		}
		err = idx.IndexTorrent(t)
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	_indexer *indexer.Indexer
	_ctx     *cli.Context
)

func getIndexer() *indexer.Indexer {
	return _indexer
}

type Config struct {
	DataDir string
	DB      DBCOnf  `envPrefix:"DB_"`
	Log     LogConf `envPrefix:"LOG_"`
}

type DBCOnf struct {
	Type     string `env:"TYPE" envDefault:"sqlite"`
	Driver   string `env:"DRIVER"`
	Database string `env:"DATABASE"`
	Username string `env:"USERNAME" envDefault:"torrenti"`
	Password string `env:"PASSWORD" envDefault:"torrenti"`
	URL      string `env:"URL"`
}

type LogConf struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

func showStat(ctx *cli.Context) error {
	idx := getIndexer()
	db := idx.DB
	st := &stat{}
	err := multierr.Combine(
		db.Model(models.MetaFile{}).Count(&st.MetaCount).Error,
		db.Model(models.MetaFile{}).Select("sum(size)").Scan(&st.MetaSize).Error,
		db.Model(models.Torrent{}).Count(&st.TorrentCount).Error,
		db.Model(models.TorrentFile{}).Count(&st.FileCount).Error,
		db.Model(models.Torrent{}).Select("sum(total_size)").Scan(&st.TorrentTotalFileSize).Error,
	)
	if err != nil {
		return err
	}

	return printYaml(map[string]interface{}{
		"meta": map[string]interface{}{
			"count":           st.MetaCount,
			"total file size": humanize.Bytes(uint64(st.MetaSize)),
		},
		"torrent": map[string]interface{}{
			"count":      st.TorrentCount,
			"total size": humanize.Bytes(uint64(st.TorrentTotalFileSize)),
			"files":      st.FileCount,
		},
	})
}

type stat struct {
	MetaCount            int64
	MetaSize             int64
	TorrentCount         int64
	TorrentTotalFileSize int64
	FileCount            int64
}
