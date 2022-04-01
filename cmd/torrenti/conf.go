package main

//go:generate gomodifytags -file=conf.go -w -all -add-tags yaml -transform snakecase --skip-unexported -add-options yaml=omitempty

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/rs/zerolog/log"
	"github.com/wenerme/torrenti/pkg/serve"
	"github.com/wenerme/torrenti/pkg/torrenti/util"
	"go.uber.org/multierr"
)

type Config struct {
	util.DirConf `yaml:",inline"`
	PluginDir    string             `env:"PLUGIN_DIR" yaml:"plugin_dir,omitempty"`
	DB           serve.DatabaseConf `envPrefix:"DB_" yaml:"db,omitempty"`
	Log          serve.LogConf      `envPrefix:"LOG_" yaml:"log,omitempty"`
	HTTP         serve.HTTPConf     `envPrefix:"HTTP_" yaml:"http,omitempty"`
	Debug        serve.DebugConf    `envPrefix:"DEBUG_" yaml:"debug,omitempty"`
	GRPC         serve.GRPCConf     `envPrefix:"GRPC_" yaml:"grpc,omitempty"`
	Scrape       ScrapeConf         `envPrefix:"SCRAPE_" yaml:"scrape,omitempty"`

	Torrent TorrentConf `envPrefix:"TORRENT_" yaml:"torrent,omitempty"`
	Sub     SubConf     `envPrefix:"SUB_" yaml:"sub,omitempty"`
}

func (conf *Config) defaults() {
	if conf.PluginDir == "" {
		self := os.Args[0]
		conf.PluginDir = filepath.Join(filepath.Dir(self), "plugins")

		if strings.Contains(self, "go-build") {
			conf.PluginDir = "bin/plugins"
		}
	}

	if conf.DB.Type != "" {
		m := map[string]interface{}{}
		err := mapstructure.Decode(conf.DB, &m)
		for k, v := range m {
			switch {
			case v == "":
			case v == 0:
			default:
				continue
			}
			delete(m, k)
		}
		err = multierr.Combine(
			err,
			mapstructure.Decode(m, &conf.Torrent.DB),
			mapstructure.Decode(m, &conf.Sub.DB),
			mapstructure.Decode(m, &conf.Scrape.Store.DB),
		)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}

	if conf.Torrent.DB.Type == "" {
		defaultTo(&conf.Sub.DB, serve.DatabaseConf{
			Type:     "sqlite",
			Database: "$DATA_DIR/torrenti.sqlite",
			Attributes: map[string]string{
				"cache": "shared",
			},
		})
	}
	if conf.Sub.DB.Type == "" {
		defaultTo(&conf.Sub.DB, serve.DatabaseConf{
			Type:     "sqlite",
			Database: "$DATA_DIR/subi.sqlite",
			Attributes: map[string]string{
				"cache": "shared",
			},
		})
	}
	if conf.Scrape.Store.DB.Type == "" {
		defaultTo(&conf.Scrape.Store.DB, serve.DatabaseConf{
			Type:     "sqlite",
			Database: "$CACHE_DIR/scraper-store.db",
			GORM: serve.GORMConf{
				Log: serve.SQLLogConf{
					IgnoreNotFound: true,
				},
			},
			Attributes: map[string]string{
				"cache": "shared",
			},
		})
	}
}

func defaultTo(a interface{}, def interface{}) {
	m := map[string]interface{}{}
	err := mapstructure.Decode(def, &m)
	err = multierr.Combine(err, mapstructure.Decode(a, &m))
	for k, v := range m {
		del := false
		switch vv := v.(type) {
		case string:
			del = len(vv) == 0
		case int:
			del = vv == 0
		case map[string]interface{}:
			del = len(vv) == 0
		case []interface{}:
			del = len(vv) == 0
		}

		if del {
			delete(m, k)
		}
	}
	err = multierr.Combine(err, mapstructure.WeakDecode(m, a))
	if err != nil {
		log.Fatal().Err(err).Msg("apply default")
	}
}

type TorrentConf struct {
	DB serve.DatabaseConf `envPrefix:"DB_" yaml:"db,omitempty"`
}
type SubConf struct {
	DB serve.DatabaseConf `envPrefix:"DB_" yaml:"db,omitempty"`
}

type ScrapeStoreConf struct {
	DB serve.DatabaseConf `envPrefix:"STORE_DB_" yaml:"db,omitempty"`
}
type ScrapeConf struct {
	Debug ScrapeDebugConf `envPrefix:"DEBUG_" yaml:"debug,omitempty"`
	Store ScrapeStoreConf `envPrefix:"STORE_" yaml:"store,omitempty"`
}
type ScrapeDebugConf struct {
	Addr string `env:"ADDR" envDefault:"127.0.0.1:7676" yaml:"addr,omitempty"`
}
