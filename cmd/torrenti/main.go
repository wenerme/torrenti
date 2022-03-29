package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/wenerme/torrenti/pkg/subi"

	"github.com/caarlos0/env/v6"
	"github.com/dustin/go-humanize"
	cli "github.com/urfave/cli/v2"
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
				Name:   "scrape",
				Usage:  "scrape web pages",
				Action: runScrape,
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "fatal",
						Usage: "stop when error occurs",
					},
					&cli.BoolFlag{
						Name:  "seed",
						Usage: "only scrape given url",
					},
					&cli.BoolFlag{
						Name:  "queue",
						Usage: "scrape pending queued url, not seed",
					},
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

var _conf = &Config{
	DB: DatabaseConf{
		Type: "sqlite",
	},
	Debug: DebugConf{
		ListenConf: util.ListenConf{
			Port: 9090,
		},
	},
	Web: WebConf{
		ListenConf: util.ListenConf{
			Port: 18080,
		},
	},
	GRPC: GRPCConf{
		ListenConf: util.ListenConf{
			Port: 18443,
		},
		Gateway: GRPCGateway{
			Prefix: "/api",
		},
	},
	Torrent: TorrentConf{
		DB: DatabaseConf{
			Type: "sqlite",
		},
	},
	Sub: SubConf{
		DB: DatabaseConf{
			Type: "sqlite",
		},
	},
}

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

	if conf.DB.Type == "sqlite" {
		conf.DB.Database = filepath.Join(conf.DataDir, "torrenti.sqlite")
	}

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

func addTorrent(ctx *cli.Context) error {
	idx := getTorrentIndexer()
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
		_, err = idx.IndexTorrent(ctx.Context, t)
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	_torrenti     *torrenti.Indexer
	_torrentiOnce = new(sync.Once)
	_subiOnce     = new(sync.Once)
	_subi         *subi.Indexer
	_ctx          *cli.Context
)

func getTorrentIndexer() *torrenti.Indexer {
	_torrentiOnce.Do(_initIndexer)
	return _torrenti
}

func getSubIndexer() *subi.Indexer {
	_subiOnce.Do(_initSubIndexer)
	return _subi
}

func _initIndexer() {
	conf := _conf
	db, gdb, err := newDB(&conf.DB)
	_ = db

	_torrenti, err = torrenti.NewIndexer(torrenti.NewIndexerOptions{DB: gdb})
	if err != nil {
		panic(err)
	}
}

func _initSubIndexer() {
	conf := _conf
	dc := &DatabaseConf{
		Type:     "sqlite",
		Database: filepath.Join(conf.DataDir, "subi.sqlite"),
	}
	db, gdb, err := newDB(dc)
	_ = db

	_subi, err = subi.NewIndexer(subi.NewIndexerOptions{DB: gdb})
	if err != nil {
		panic(err)
	}
}

type Config struct {
	util.DirConf `yaml:",inline"`
	PluginDir    string       `env:"PLUGIN_DIR"`
	DB           DatabaseConf `envPrefix:"DB_"`
	Log          LogConf      `envPrefix:"LOG_"`
	GRPC         GRPCConf     `envPrefix:"GRPC_"`
	Web          WebConf      `envPrefix:"WEB_"`
	Debug        DebugConf    `envPrefix:"DEBUG_"`
	Scrape       ScrapeConf   `envPrefix:"SCRAPE_"`

	Torrent TorrentConf `envPrefix:"TORRENT_"`
	Sub     SubConf     `envPrefix:"SUB_"`
}

type DatabaseConf struct {
	Type     string     `env:"TYPE"`
	Driver   string     `env:"DRIVER"`
	Database string     `env:"DATABASE"`
	Username string     `env:"USERNAME"`
	Password string     `env:"PASSWORD"`
	Host     string     `env:"HOST"`
	Port     string     `env:"PORT"`
	Schema   string     `env:"SCHEMA"`
	DSN      string     `env:"DSN"`
	Log      SQLLogConf `envPrefix:"LOG_"`

	DriverOptions DatabaseDriverOptions `envPrefix:"DRIVER_"`
	Attributes    map[string]string     `envPrefix:"ATTR_"` // ConnectionAttributes
}

type DatabaseDriverOptions struct {
	MaxIdleConnections int            `env:"MAX_IDLE_CONNS"  envDefault:"20"`
	MaxOpenConnections int            `env:"MAX_OPEN_CONNS"  envDefault:"100"`
	ConnMaxIdleTime    time.Duration  `env:"MAX_IDLE_TIME"  envDefault:"10m"`
	ConnMaxLifetime    *time.Duration `env:"MAX_LIVE_TIME"`
}

type SQLLogConf struct {
	SlowThreshold  time.Duration `env:"SLOW_THRESHOLD"`
	IgnoreNotFound bool          `env:"IGNORE_NOT_FOUND"`
	Debug          bool          `env:"DEBUG"`
}

type LogConf struct {
	Level string `env:"LEVEL" envDefault:"info"`
}

func showStat(ctx *cli.Context) error {
	idx := getTorrentIndexer()
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

type GRPCConf struct {
	util.ListenConf `yaml:",inline"`
	Enabled         bool        `env:"ENABLED" envDefault:"true"`
	Gateway         GRPCGateway `envPrefix:"GATEWAY_"`
}
type GRPCGateway struct {
	util.ListenConf `yaml:",inline"`
	Enabled         bool   `env:"ENABLED" envDefault:"true"`
	Prefix          string `env:"PREFIX"`
}

type WebConf struct {
	util.ListenConf `yaml:",inline"`
}
type DebugConf struct {
	util.ListenConf `yaml:",inline"`
	Enabled         bool `env:"ENABLED" envDefault:"true"`
}
type TorrentConf struct {
	DB DatabaseConf `envPrefix:"DB_"`
}
type SubConf struct {
	DB DatabaseConf `envPrefix:"DB_"`
}

type ScrapeConf struct {
	Debug ScrapeDebugConf `envPrefix:"DEBUG_"`
}
type ScrapeDebugConf struct {
	Addr string `env:"ADDR" envDefault:"127.0.0.1:7676"`
}
