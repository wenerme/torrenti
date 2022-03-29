package main

import (
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/wenerme/torrenti/pkg/scrape"
	"gorm.io/gorm"

	"github.com/wenerme/torrenti/pkg/plugin"

	"github.com/wenerme/torrenti/pkg/subi"

	"github.com/caarlos0/env/v6"
	"github.com/dustin/go-humanize"
	cli "github.com/urfave/cli/v2"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"gopkg.in/yaml.v3"

	"github.com/wenerme/torrenti/pkg/torrenti/models"
	"go.uber.org/multierr"

	_ "github.com/glebarez/go-sqlite"
	_ "github.com/jackc/pgx/v4"
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
				Name:   "migration",
				Action: runMigration,
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
						Name:  "pull",
						Usage: "pull pending queued url, not seed",
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
	DB: DatabaseConf{},
	Debug: DebugConf{
		ListenConf: util.ListenConf{
			Port: 9090,
		},
	},
	HTTP: HTTPConf{
		ListenConf: util.ListenConf{
			Port: 18080,
		},
	},
	GRPC: GRPCConf{
		ListenConf: util.ListenConf{
			Port: 18443,
		},
		Gateway: GRPCGatewayConf{
			Prefix: "/api",
		},
	},
	Torrent: TorrentConf{
		DB: DatabaseConf{},
	},
	Sub: SubConf{
		DB: DatabaseConf{},
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

	if cfgFile == "" {
		if _, err := os.Stat("./config.yaml"); err == nil {
			cfgFile = "./config.yaml"
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

	conf.InitDirConf(Name)
	conf.defaults()

	log.Debug().Str("config_dir", conf.ConfigDir).Msg("conf")
	log.Debug().Str("data_dir", conf.DataDir).Msg("conf")
	log.Debug().Str("cache_dir", conf.CacheDir).Msg("conf")
	log.Debug().Str("plugin_dir", conf.PluginDir).Msg("conf")
	log.Debug().Str("log_level", conf.Log.Level).Msg("conf")

	if err := multierr.Combine(
		os.Setenv("DATA_DIR", conf.DataDir),
		os.Setenv("CACHE_DIR", conf.CacheDir),
		os.Setenv("CONFIG_DIR", conf.ConfigDir),
	); err != nil {
		log.Fatal().Err(err).Msg("set dir env")
	}

	err := os.MkdirAll(conf.DataDir, 0o755)
	if err != nil {
		return err
	}
	_ctx = ctx
	err = plugin.LoadPlugins(plugin.LoadPluginOptions{
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
	_, gdb, err := newDB(&conf.Torrent.DB)
	if err != nil {
		panic(err)
	}

	_torrenti, err = torrenti.NewIndexer(torrenti.NewIndexerOptions{DB: gdb})
	if err != nil {
		panic(err)
	}
}

func _initSubIndexer() {
	conf := _conf
	_, gdb, err := newDB(&conf.Sub.DB)
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	_subi, err = subi.NewIndexer(subi.NewIndexerOptions{DB: gdb})
	if err != nil {
		panic(err)
	}
}

func runMigration(ctx *cli.Context) (err error) {
	getTorrentIndexer()
	getSubIndexer()

	var sdb *gorm.DB
	_, sdb, err = newDB(&_conf.Scrape.Store.DB)
	if err != nil {
		return errors.Wrap(err, "new scraper store")
	}
	sc := &scrape.Context{
		DB: sdb,
	}
	sc.Seed, _ = url.Parse("https://wener.me")
	return sc.Init()
}

func showStat(ctx *cli.Context) error {
	idx := getTorrentIndexer()
	db := idx.DB
	st := &stat{}
	ts, err := idx.Stat(ctx.Context)
	err = multierr.Combine(
		err,
		db.Model(models.Torrent{}).Select("coalesce(max(total_file_size),0)").Scan(&st.MaxTorrentSize).Error,
		db.Model(models.Torrent{}).Where(models.Torrent{TotalFileSize: st.MaxTorrentSize}).Select("name").Scan(&st.MaxTorrentName).Error,
	)
	if err != nil {
		return err
	}
	return printYaml(map[string]interface{}{
		"meta": map[string]interface{}{
			"count":           ts.MetaCount,
			"total file size": humanize.Bytes(uint64(ts.MetaSize)),
		},
		"torrent": map[string]interface{}{
			"count":            ts.TorrentCount,
			"total size":       humanize.Bytes(uint64(ts.TorrentFileTotalSize)),
			"files":            ts.TorrentFileCount,
			"max torrent size": humanize.Bytes(uint64(st.MaxTorrentSize)),
			"max torrent name": st.MaxTorrentName,
		},
	})
}

type stat struct {
	MaxTorrentSize int64
	MaxTorrentName string
}
