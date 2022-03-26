package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/caarlos0/env/v6"
	"github.com/dustin/go-humanize"
	"github.com/gocolly/colly"
	cli "github.com/urfave/cli/v2"
	"github.com/wenerme/torrenti/pkg/torrenti/scraper"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"gopkg.in/yaml.v3"

	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"go.uber.org/multierr"

	_ "github.com/glebarez/go-sqlite"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/wenerme/torrenti/pkg/torrenti"
)

const Name = "torrenti"

func main() {
	app := &cli.App{
		Name:  Name,
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
				Name:   "serve",
				Action: runServer,
			},
			{
				Name: "version",
				Action: func(c *cli.Context) error {
					info := util.ReadBuildInfo()
					fmt.Println(Name, info.String())
					return nil
				},
			},
			{
				Name: "torrent",
				Subcommands: cli.Commands{
					{
						Name:   "add",
						Usage:  "add to index",
						Action: addTorrent,
					},
					{
						Name:   "scrape",
						Usage:  "scrape web pages",
						Action: scrapeTorrent,
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:  "fatal",
								Usage: "stop when error occurs",
							},
						},
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

var _conf = &Config{}

func setup(ctx *cli.Context) error {
	conf := _conf

	{
		l := zerolog.InfoLevel
		if lvl, err := zerolog.ParseLevel(ctx.String("log-level")); err == nil {
			l = lvl
		}
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		log.Logger = zerolog.New(output).Level(l).With().Stack().Timestamp().Logger()
	}

	cfgFile := os.ExpandEnv(ctx.String("config"))
	if cfgFile == "" {
		cfgFile = filepath.Join(conf.ConfigDir, "config.yaml")
		if _, err := os.Stat(cfgFile); err != nil {
			log.Debug().Str("path", cfgFile).Msg("default config file not found")
			cfgFile = ""
		}
	}

	if cfgFile != "" {
		log.Debug().Str("path", cfgFile).Msg("use config file")

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
	if conf.PluginDir == "" {
		self := os.Args[0]
		conf.PluginDir = filepath.Join(filepath.Dir(self), "plugins")

		if strings.Contains(self, "go-build") {
			conf.PluginDir = "bin/plugins"
		}
	}
	conf.InitDirConf(Name)

	log.Debug().Str("config_dir", conf.ConfigDir).Msg("conf")
	log.Debug().Str("data_dir", conf.DataDir).Msg("conf")
	log.Debug().Str("cache_dir", conf.CacheDir).Msg("conf")
	log.Debug().Str("plugin_dir", conf.PluginDir).Msg("conf")
	log.Debug().Str("log_level", conf.Log.Level).Msg("conf")

	err := os.MkdirAll(conf.DataDir, 0o755)
	if err != nil {
		return err
	}
	_ctx = ctx
	err = torrenti.LoadPlugins(torrenti.LoadPluginOptions{
		Dir: conf.PluginDir,
	})
	if err != nil {
		return err
	}

	return err
}

func scrapeTorrent(ctx *cli.Context) error {
	cache := filepath.Join(_conf.CacheDir, "web")
	if err := os.MkdirAll(cache, 0o755); err != nil {
		log.Fatal().Err(err).Msg("make cache dir")
	}
	log.Debug().Str("cache_dir", cache).Msg("cache dir")
	if ctx.NArg() != 1 {
		log.Fatal().Msgf("must scrap only one url got %v", ctx.NArg())
	}
	first := ctx.Args().First()
	u, err := url.Parse(first)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid url")
	}
	opts := &scraper.ScrapeOptions{
		Fatal: !ctx.Bool("fatal"),
	}
	ct := context.Background()
	ct = torrenti.IndexerContextKey.WithValue(ct, getIndexer())
	ct = util.DirConfContextKey.WithValue(ct, &_conf.DirConf)
	ct = scraper.OptionContextKey.WithValue(ct, opts)

	c := colly.NewCollector(
		colly.CacheDir(cache),
		scraper.InitCollector(ct, u),
	)

	scraper.SetupCollector(ct, u)(c)
	return c.Visit(first)
}

func addTorrent(ctx *cli.Context) error {
	idx := getIndexer()
	for _, v := range ctx.Args().Slice() {
		log.Info().Str("torrent", v).Msg("add torrent")
		t, err := torrenti.ParseTorrent(v)
		if err != nil {
			return err
		}
		err = t.Load()
		if err != nil {
			return err
		}
		_, err = idx.IndexTorrent(t)
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	_indexer     *torrenti.Indexer
	_indexerOnce = new(sync.Once)
	_ctx         *cli.Context
)

func getIndexer() *torrenti.Indexer {
	_indexerOnce.Do(_initIndexer)
	return _indexer
}

func _initIndexer() {
	conf := _conf
	var err error
	var gdb *gorm.DB
	var db *sql.DB
	var dsn string
	var driver string
	switch conf.DB.Type {
	case "sqlite":
		dbFile := conf.DB.Database
		if conf.DB.Database == "" {
			dbFile = filepath.Join(conf.DataDir, "torrenti.sqlite")
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
		//if err = gdb.Exec("PRAGMA page_size = ?", 128*1024).Error; err != nil {
		//	return err
		//}
	default:
		err = errors.New("unsupported db type: " + conf.DB.Type)
	}
	if err != nil {
		panic(err)
	}

	_indexer, err = torrenti.NewIndexer(torrenti.NewIndexerOptions{DB: gdb})
	if err != nil {
		panic(err)
	}
}

type Config struct {
	util.DirConf `yaml:",inline"`
	PluginDir    string  `env:"PLUGIN_DIR"`
	DB           DBCOnf  `envPrefix:"DB_"`
	Log          LogConf `envPrefix:"LOG_"`
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

func runServer(ctx *cli.Context) error {
	return nil
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
		db.Model(models.Torrent{}).Select("sum(total_file_size)").Scan(&st.TorrentTotalFileSize).Error,
		db.Model(models.Torrent{}).Select("max(total_file_size)").Scan(&st.MaxTorrentSize).Error,
	)
	if err != nil {
		return err
	}
	err = multierr.Combine(
		db.Model(models.Torrent{}).Where(models.Torrent{TotalFileSize: st.MaxTorrentSize}).Select("name").Scan(&st.MaxTorrentName).Error,
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
			"count":            st.TorrentCount,
			"total size":       humanize.Bytes(uint64(st.TorrentTotalFileSize)),
			"files":            st.FileCount,
			"max torrent size": humanize.Bytes(uint64(st.MaxTorrentSize)),
			"max torrent name": st.MaxTorrentName,
		},
	})
}

type stat struct {
	MetaCount            int64
	MetaSize             int64
	TorrentCount         int64
	MaxTorrentSize       int64
	MaxTorrentName       string
	TorrentTotalFileSize int64
	FileCount            int64
}
